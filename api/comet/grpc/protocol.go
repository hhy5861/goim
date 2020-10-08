package grpc

import (
	"encoding/json"
	"errors"
	"github.com/hhy5861/goim/pkg/bufio"
	"github.com/hhy5861/goim/pkg/bytes"
	"github.com/hhy5861/goim/pkg/websocket"
	"reflect"
	"strconv"
)

const (
	// MaxBodySize max proto body size
	MaxBodySize = int32(1 << 12)
)

const (
	// size
	_packSize      = 4
	_headerSize    = 2
	_verSize       = 2
	_opSize        = 4
	_seqSize       = 4
	_rawHeaderSize = _packSize + _headerSize + _verSize + _opSize + _seqSize
	// offset
)

var (
	// ErrProtoPackLen proto packet len error
	ErrProtoPackLen = errors.New("default server codec pack length error")
)

// WriteTo write a proto to bytes writer.
func (p *Proto) WriteTo(b *bytes.Writer) {

	if p.Body != nil {
		body, err := json.Marshal(p.Body)
		if err != nil {
			return
		}

		b.Write(body)
	}
}

// ReadTCP read a proto from TCP reader.
func (p *Proto) ReadTCP(rr *bufio.Reader) (err error) {
	var (
		buf []byte
	)

	buf, err = rr.Pop(_rawHeaderSize)
	if err != nil {
		return
	}

	if len(buf) > 0 {
		var msg *ClientMsg
		if err := json.Unmarshal(buf, msg); err != nil {
			return err
		}

		p.Body = msg
	} else {
		p.Body = nil
	}

	return
}

// WriteTCP write a proto to TCP writer.
func (p *Proto) WriteTCP(wr *bufio.Writer) (err error) {
	if p.Body != nil {
		body, err := json.Marshal(p.Body)
		if err != nil {
			return
		}

		_, err = wr.Write(body)
	}

	return
}

// WriteTCPHeart write TCP heartbeat with room online.
func (p *Proto) WriteTCPHeart(wr *bufio.Writer, online int32) (err error) {
	_, err = wr.Write([]byte(strconv.Itoa(int(online))))
	return
}

// ReadWebsocket read a proto from websocket connection.
func (p *Proto) ReadWebsocket(ws *websocket.Conn) (err error) {
	var (
		buf []byte
	)

	if _, buf, err = ws.ReadMessage(); err != nil {
		return
	}
	if len(buf) < _rawHeaderSize {
		return ErrProtoPackLen
	}

	var msg *ClientMsg
	if err := json.Unmarshal(buf, msg); err != nil {
		return err
	}

	if p.IsNil(msg) == false {
		p.Body = msg
	} else {
		p.Body = nil
	}

	return
}

// WriteWebsocket write a proto to websocket connection.
func (p *Proto) WriteWebsocket(ws *websocket.Conn) (err error) {
	if p.Body != nil {
		body, err := json.Marshal(p.Body)
		if err != nil {
			return err
		}

		err = ws.WriteMessage(websocket.TextMessage, body)
	}

	return
}

// WriteWebsocketHeart write websocket heartbeat with room online.
func (p *Proto) WriteWebsocketHeart(wr *websocket.Conn, online int32) (err error) {
	err = wr.WriteMessage(websocket.PingMessage, []byte(strconv.Itoa(int(online))))
	return
}

func (p *Proto) IsNil(v interface{}) bool {
	return v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil())
}
