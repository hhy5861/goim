package common

import "time"

type (
	MsgGetOpts struct {
		User            string     `json:"user,omitempty"`
		Topic           string     `json:"topic,omitempty"`
		IfModifiedSince *time.Time `json:"ims,omitempty"`
		SinceId         int        `json:"since,omitempty"`
		BeforeId        int        `json:"before,omitempty"`
		Limit           int        `json:"limit,omitempty"`
	}

	MsgSetSub struct {
		User string `json:"user,omitempty"`
		Mode string `json:"mode,omitempty"`
	}

	MsgSetQuery struct {
		Sub  *MsgSetSub `json:"sub,omitempty"`
		Tags []string   `json:"tags,omitempty"`
	}

	MsgClientHi struct {
		Id        string `json:"id,omitempty"`
		UserAgent string `json:"ua,omitempty"`
		Version   string `json:"ver,omitempty"`
		DeviceID  string `json:"dev,omitempty"`
		Lang      string `json:"lang,omitempty"`
		Platform  string `json:"platf,omitempty"`
	}

	MsgClientLogin struct {
		Id     string `json:"id,omitempty"`
		Scheme string `json:"scheme,omitempty"`
		Secret []byte `json:"secret"`
	}

	MsgGetQuery struct {
		What string      `json:"what"`
		Sub  *MsgGetOpts `json:"sub,omitempty"`
		Data *MsgGetOpts `json:"data,omitempty"`
	}

	MsgClientSub struct {
		Id         string       `json:"id,omitempty"`
		Topic      string       `json:"topic"`
		Background bool         `json:"bkg,omitempty"`
		Set        *MsgSetQuery `json:"set,omitempty"`
		Get        *MsgGetQuery `json:"get,omitempty"`
	}

	MsgClientLeave struct {
		Id    string `json:"id,omitempty"`
		Topic string `json:"topic"`
		Unsub bool   `json:"unsub,omitempty"`
	}

	MsgClientPub struct {
		Id      string                 `json:"id,omitempty"`
		Topic   string                 `json:"topic"`
		NoEcho  bool                   `json:"noecho,omitempty"`
		Head    map[string]interface{} `json:"head,omitempty"`
		Content interface{}            `json:"content"`
	}

	MsgClientGet struct {
		Id    string `json:"id,omitempty"`
		Topic string `json:"topic"`
		MsgGetQuery
	}

	ClientComMessage struct {
		Hi        *MsgClientHi    `json:"hi"`
		Login     *MsgClientLogin `json:"login"`
		Sub       *MsgClientSub   `json:"sub"`
		Leave     *MsgClientLeave `json:"leave"`
		Pub       *MsgClientPub   `json:"pub"`
		Get       *MsgClientGet   `json:"get"`
		Id        string
		Topic     string
		AsUser    string
		AuthLvl   int
		Timestamp time.Time
	}
)

type (
	MsgServerCtrl struct {
		Id        string      `json:"id,omitempty"`
		Topic     string      `json:"topic,omitempty"`
		Params    interface{} `json:"params,omitempty"`
		Code      int         `json:"code"`
		Text      string      `json:"text,omitempty"`
		Timestamp time.Time   `json:"ts"`
	}

	MsgServerData struct {
		Topic     string                 `json:"topic"`
		From      string                 `json:"from,omitempty"`
		Timestamp time.Time              `json:"ts"`
		DeletedAt *time.Time             `json:"deleted,omitempty"`
		SeqId     int                    `json:"seq"`
		Head      map[string]interface{} `json:"head,omitempty"`
		Content   interface{}            `json:"content"`
		Internal  string                 `json:"internal,omitempty"`
	}

	ServerComMessage struct {
		Ctrl      *MsgServerCtrl `json:"ctrl,omitempty"`
		Data      *MsgServerData `json:"data,omitempty"`
		Id        string
		Timestamp time.Time
		AsUser    string
		SkipSid   string
	}
)
