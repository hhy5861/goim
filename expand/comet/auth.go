package comet

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
)

const (
	apikeyVersion   = 1
	apikeyAppID     = 4
	apikeySequence  = 2
	apikeyWho       = 1
	apikeySignature = 16
	apikeyLength    = apikeyVersion + apikeyAppID + apikeySequence + apikeyWho + apikeySignature

	ApiKeyParamName = "apikey"
)

type (
	Auth struct {
		keySalt string
	}
)

func NewAuth(salt string) *Auth {
	return &Auth{
		keySalt: salt,
	}
}

func (svc *Auth) CheckApiKey(apikey string) (isValid bool) {
	if declen := base64.URLEncoding.DecodedLen(len(apikey)); declen != apikeyLength {
		return
	}

	data, err := base64.URLEncoding.DecodeString(apikey)
	if err != nil {
		log.Println("failed to decode.base64 appid ", err)
		return
	}

	if data[0] != 1 {
		log.Println("unknown appid signature algorithm ", data[0])
		return
	}

	keySalt, err := base64.StdEncoding.DecodeString(svc.keySalt)
	if err != nil {
		log.Println("api key salt marshal", err)
		return
	}

	hasher := hmac.New(md5.New, keySalt)
	hasher.Write(data[:apikeyVersion+apikeyAppID+apikeySequence+apikeyWho])
	check := hasher.Sum(nil)

	if !bytes.Equal(data[apikeyVersion+apikeyAppID+apikeySequence+apikeyWho:], check) {
		log.Println("invalid api key signature")
		return
	}

	return true
}

func (svc *Auth) BuildApiPath(serviceName, path string) string {
	path = strings.TrimLeft(path, "/")

	if serviceName != "" && serviceName != "/" {
		return fmt.Sprintf("%s/%s", strings.TrimLeft(serviceName, "/"), path)
	}

	if serviceName == "/" {
		return fmt.Sprintf("%s%s", serviceName, path)
	}

	return fmt.Sprintf("/%s", path)
}
