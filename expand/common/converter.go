package common

import (
	"encoding/json"
	pbx "github.com/hhy5861/goim/api/comet/grpc"
	"log"
	"time"
)

func PbCliDeserialize(pkt *pbx.ClientMsg) *ClientComMessage {
	var msg ClientComMessage
	if hi := pkt.GetHi(); hi != nil {
		msg.Hi = &MsgClientHi{
			Id:        hi.GetId(),
			UserAgent: hi.GetUserAgent(),
			Version:   hi.GetVer(),
			DeviceID:  hi.GetDeviceId(),
			Platform:  hi.GetPlatform(),
			Lang:      hi.GetLang(),
		}
	} else if login := pkt.GetLogin(); login != nil {
		msg.Login = &MsgClientLogin{
			Id:     login.GetId(),
			Scheme: login.GetScheme(),
			Secret: login.GetSecret(),
		}
	} else if sub := pkt.GetSub(); sub != nil {
		msg.Sub = &MsgClientSub{
			Id:         sub.GetId(),
			Topic:      sub.GetTopic(),
			Background: sub.GetBackground(),
			Get:        pbGetQueryDeserialize(sub.GetGetQuery()),
			Set:        pbSetQueryDeserialize(sub.GetSetQuery()),
		}
	} else if leave := pkt.GetLeave(); leave != nil {
		msg.Leave = &MsgClientLeave{
			Id:    leave.GetId(),
			Topic: leave.GetTopic(),
			Unsub: leave.GetUnsub(),
		}
	} else if pub := pkt.GetPub(); pub != nil {
		msg.Pub = &MsgClientPub{
			Id:      pub.GetId(),
			Topic:   pub.GetTopic(),
			NoEcho:  pub.GetNoEcho(),
			Head:    byteMapToInterfaceMap(pub.GetHead()),
			Content: bytesToInterface(pub.GetContent()),
		}
	} else if get := pkt.GetGet(); get != nil {
		msg.Get = &MsgClientGet{
			Id:          get.GetId(),
			Topic:       get.GetTopic(),
			MsgGetQuery: *pbGetQueryDeserialize(get.GetQuery()),
		}
	}

	msg.AsUser = pkt.GetOnBehalfOf()
	return &msg
}

// Convert ClientComMessage to pbx.ClientMsg
func PbCliSerialize(msg *ClientComMessage) *pbx.ClientMsg {
	var pkt *pbx.ClientMsg

	switch {
	case msg.Hi != nil:
		pkt = &pbx.ClientMsg{
			Hi: &pbx.ClientHi{
				Id:        msg.Hi.Id,
				UserAgent: msg.Hi.UserAgent,
				Ver:       msg.Hi.Version,
				DeviceId:  msg.Hi.DeviceID,
				Platform:  msg.Hi.Platform,
				Lang:      msg.Hi.Lang,
			}}
	case msg.Login != nil:
		pkt = &pbx.ClientMsg{
			Login: &pbx.ClientLogin{
				Id:     msg.Login.Id,
				Scheme: msg.Login.Scheme,
				Secret: msg.Login.Secret,
			}}
	case msg.Sub != nil:
		pkt = &pbx.ClientMsg{Sub: &pbx.ClientSub{
			Id:         msg.Sub.Id,
			Topic:      msg.Sub.Topic,
			Background: msg.Sub.Background,
			SetQuery:   pbSetQuerySerialize(msg.Sub.Set),
			GetQuery:   pbGetQuerySerialize(msg.Sub.Get)}}
	case msg.Leave != nil:
		pkt = &pbx.ClientMsg{Leave: &pbx.ClientLeave{
			Id:    msg.Leave.Id,
			Topic: msg.Leave.Topic,
			Unsub: msg.Leave.Unsub}}
	case msg.Pub != nil:
		pkt = &pbx.ClientMsg{Pub: &pbx.ClientPub{
			Id:      msg.Pub.Id,
			Topic:   msg.Pub.Topic,
			NoEcho:  msg.Pub.NoEcho,
			Head:    interfaceMapToByteMap(msg.Pub.Head),
			Content: interfaceToBytes(msg.Pub.Content)}}
	case msg.Get != nil:
		pkt = &pbx.ClientMsg{Get: &pbx.ClientGet{
			Id:    msg.Get.Id,
			Topic: msg.Get.Topic,
			Query: pbGetQuerySerialize(&msg.Get.MsgGetQuery)}}
	}

	if pkt == nil {
		return nil
	}

	pkt.OnBehalfOf = msg.AsUser
	return pkt
}

func pbGetQueryDeserialize(in *pbx.GetQuery) *MsgGetQuery {
	msg := MsgGetQuery{}

	if in != nil {
		msg.What = in.GetWhat()

		if sub := in.GetSub(); sub != nil {
			msg.Sub = &MsgGetOpts{
				IfModifiedSince: int64ToTime(sub.GetIfModifiedSince()),
				Limit:           int(sub.GetLimit()),
			}
		}
		if data := in.GetData(); data != nil {
			msg.Data = &MsgGetOpts{
				BeforeId: int(data.GetBeforeId()),
				SinceId:  int(data.GetSinceId()),
				Limit:    int(data.GetLimit()),
			}
		}
	}

	return &msg
}

func pbSetQueryDeserialize(in *pbx.SetQuery) *MsgSetQuery {
	var msg *MsgSetQuery

	if in != nil {
		if sub := in.GetSub(); sub != nil {
			user := sub.GetUserId()
			mode := sub.GetMode()

			if user != "" || mode != "" {
				if msg == nil {
					msg = &MsgSetQuery{}
				}

				msg.Sub = &MsgSetSub{
					User: sub.GetUserId(),
					Mode: sub.GetMode(),
				}
			}
		}
	}

	return msg
}

func byteMapToInterfaceMap(in map[string][]byte) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for key, raw := range in {
		if val := bytesToInterface(raw); val != nil {
			out[key] = val
		}
	}
	return out
}

func interfaceMapToByteMap(in map[string]interface{}) map[string][]byte {
	out := make(map[string][]byte, len(in))
	for key, val := range in {
		if val != nil {
			out[key], _ = json.Marshal(val)
		}
	}
	return out
}

func bytesToInterface(in []byte) interface{} {
	var out interface{}
	if len(in) > 0 {
		err := json.Unmarshal(in, &out)
		if err != nil {
			log.Println("pbx: failed to parse bytes", string(in), err)
		}
	}
	return out
}

func interfaceToBytes(in interface{}) []byte {
	if in != nil {
		out, _ := json.Marshal(in)
		return out
	}
	return nil
}

func int64ToTime(ts int64) *time.Time {
	if ts > 0 {
		res := time.Unix(ts/1000, ts%1000).UTC()
		return &res
	}
	return nil
}

func timeToInt64(ts *time.Time) int64 {
	if ts != nil {
		return ts.UnixNano() / int64(time.Millisecond)
	}
	return 0
}

func pbSetQuerySerialize(in *MsgSetQuery) *pbx.SetQuery {
	if in == nil {
		return nil
	}

	var out *pbx.SetQuery

	if in.Sub != nil {
		out.Sub = &pbx.SetSub{
			UserId: in.Sub.User,
			Mode:   in.Sub.Mode,
		}
	}

	return out
}

func pbGetQuerySerialize(in *MsgGetQuery) *pbx.GetQuery {
	if in == nil {
		return nil
	}

	out := &pbx.GetQuery{
		What: in.What,
	}

	if in.Sub != nil {
		out.Sub = &pbx.GetOpts{
			IfModifiedSince: timeToInt64(in.Sub.IfModifiedSince),
			User:            in.Sub.User,
			Topic:           in.Sub.Topic,
			Limit:           int32(in.Sub.Limit)}
	}

	if in.Data != nil {
		out.Data = &pbx.GetOpts{
			BeforeId: int32(in.Data.BeforeId),
			SinceId:  int32(in.Data.SinceId),
			Limit:    int32(in.Data.Limit)}
	}

	return out
}
