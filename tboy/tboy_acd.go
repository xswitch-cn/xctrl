package tboy

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/xswitch-cn/proto/go/proto/xctrl"
	"github.com/xswitch-cn/xctrl/ctrl"
	"github.com/xswitch-cn/xctrl/ctrl/nats"
)

type TBoyACD struct {
	*TBoy
}

// NewTBoy .
func NewACD(node_uuid string, domain string, fn ...OptionFn) *TBoyACD {
	boy := &TBoyACD{
		TBoy: &TBoy{
			Options: &Options{},
		},
	}
	boy.Init()
	boy.SetUUID(node_uuid)
	boy.SetDomain(domain)
	for _, o := range fn {
		o(boy.Options)
	}

	log.Infof("new boy created uuid=%s domain=%s", boy.NodeUUID, boy.Domain)

	return boy
}

func (boy *TBoyACD) AddChannel(key string, channel *FakeChannel) {
	Channels[key] = channel
}

func (boy *TBoyACD) Event(msg *ctrl.Message, natsEvent nats.Event) {
	topic := natsEvent.Topic()
	reply := natsEvent.Reply()
	log.Infof("Handle %s %s", topic, msg.Method)

	ctx := context.Background()

	switch msg.Method {

	case "XNode.Transfer":
		boy.Transfer(ctx, msg, reply)
	case "XNode.Play":
		boy.Play(ctx, msg, reply)
	case "XNode.Dial":
		boy.Dial(ctx, msg, reply)
	default:
		// handle messages in TBoy
		boy.TBoy.Event(msg, natsEvent)
	}

}

func (boy *TBoyACD) Play(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.Request

	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		return
	}
	channel, ok := Channels[request.Uuid]
	if !ok {
		res := xctrl.Response{
			Code:     http.StatusBadRequest,
			Message:  "can not locate session",
			NodeUuid: boy.NodeUUID,
			Uuid:     request.Uuid,
		}
		res_bytes, _ := json.MarshalIndent(res, "", "  ")
		raw := json.RawMessage(res_bytes)

		result := ctrl.Result{
			Version: "2.0",
			ID:      msg.ID,
			Result:  &raw,
		}

		result_bytes, _ := json.MarshalIndent(result, "", "  ")
		if reply == "" {
			reply = "cn.xswitch.ctrl." + request.CtrlUuid
		}
		ctrl.Publish(reply, result_bytes)
		return
	}
	if channel.Context == nil {
		channel.Context, channel.Cancel = context.WithCancel(context.Background())
	}

	go func() {
		t := time.NewTimer(3 * time.Second)
		select {
		case <-channel.Context.Done():
			return
		case <-t.C:
			boy.OK(ctx, msg, reply)
		}
	}()

}

func (boy *TBoyACD) Transfer(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.SetVarRequest

	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		return
	}

	channel, ok := Channels[request.Uuid]

	var res xctrl.Response
	if ok {
		channel.Lock.Lock()
		for k, v := range request.Data {
			channel.Data.Params[k] = v
		}
		channel.Lock.Unlock()
		res = xctrl.Response{
			Code:     200,
			Message:  "OK",
			NodeUuid: boy.NodeUUID,
			Uuid:     request.Uuid,
		}
	} else {
		res = xctrl.Response{
			Code:     http.StatusBadRequest,
			Message:  "can not locate session",
			NodeUuid: boy.NodeUUID,
			Uuid:     request.Uuid,
		}
	}

	res_bytes, _ := json.MarshalIndent(res, "", "  ")
	raw := json.RawMessage(res_bytes)

	result := ctrl.Result{
		Version: "2.0",
		ID:      msg.ID,
		Result:  &raw,
	}

	result_bytes, _ := json.MarshalIndent(result, "", "  ")
	if reply == "" {
		reply = "cn.xswitch.ctrl." + request.CtrlUuid
	}
	ctrl.Publish(reply, result_bytes)
	channel.Data.CidName = channel.Data.Params["xcc_origin_cid_number"]
	channel.Data.CidNumber = channel.Data.Params["xcc_origin_cid_number"]
	channel.Data.DestNumber = channel.Data.Params["xcc_origin_dest_number"]
	channel.Data.State = "READY"

	event_req := ctrl.Request{
		Version: "2.0",
		Method:  "Event.Channel",
		Params:  ctrl.ToRawMessage(channel.Data),
	}

	controller := "cn.xswitch.ctrl." + request.CtrlUuid
	req_str, _ := json.MarshalIndent(event_req, "", "  ")
	ctrl.Publish(controller, req_str)

	params := Dialplan{
		UUID:              request.Uuid,
		Context:           "default",
		Session:           channel.Data.Params["sip_h_X-FS-Session"],
		Domain:            channel.Data.Params["xcc_domain"],
		Direction:         "outbound",
		NodeUUID:          boy.NodeUUID,
		CallerNumber:      channel.Data.CidNumber,
		DestinationNumber: channel.Data.DestNumber,
		OriginCidNumber:   channel.Data.Params["xcc_origin_cid_number"],
		OriginDestNumber:  channel.Data.Params["xcc_origin_dest_number"],
		CallerName:        "Outbound Call",
	}

	req := ctrl.Request{
		Version: "2.0",
		ID:      ctrl.ToRawMessage(request.Uuid),
		Method:  "XCtrl.Dialplan",
		Params:  ctrl.ToRawMessage(params),
	}
	natsMsg, err := ctrl.Call("cn.xswitch.request", &req, time.Second)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info(string(natsMsg.Body))

	var response ctrl.Response

	err = json.Unmarshal(natsMsg.Body, &response)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(response.Result)
	resByte, err := json.Marshal(response.Result)
	if err != nil {
		log.Error("response result not ok")
		return
	}
	apps := make([]*App, 0)
	err = json.Unmarshal(resByte, &apps)
	if err != nil {
		log.Error(err)
		return
	}

	var xccAction string
	for _, app := range apps {
		if strings.Contains(app.Data, "xcc_action=") {
			xccAction = strings.TrimPrefix(app.Data, "xcc_action=")
			break
		}

	}

	channel.Data.State = "START"
	channel.Data.Params["xcc_action"] = xccAction
	channel.Data.Params["xcc_mark"] = "Dialing"
	channel.Data.Params["created_time"] = strconv.Itoa(int(channel.Data.CreateEpoch * 1000000))
	event_req = ctrl.Request{
		Version: "2.0",
		Method:  "Event.Channel",
		Params:  ctrl.ToRawMessage(channel.Data),
	}

	req_str, _ = json.MarshalIndent(event_req, "", "  ")
	ctrl.Publish("cn.xswitch.ctrl.app.acd", req_str)

	if !boy.Options.AcdAssign {
		time.Sleep(time.Second * 3)
		if channel.Cancel != nil {
			channel.Cancel()
		}
		controller := "cn.xswitch.ctrl.app.acd"

		channel.Data.State = "DESTROY"
		channel.Data.Cause = "NORMAL_CLEARING"
		event_req := ctrl.Request{
			Version: "2.0",
			Method:  "Event.Channel",
			Params:  ctrl.ToRawMessage(channel.Data),
		}
		req_str, _ := json.MarshalIndent(event_req, "", "  ")
		ctrl.Publish(controller, req_str)

		cdr := new(CallDetailRecord)
		err = json.Unmarshal([]byte(cdr_template), cdr)

		if err != nil {
			log.Error(err)
			return
		}

		if send_cdr {

			cdr.UUID = channel.Data.Uuid
			cdr.Leg = "a"
			cdr.StartStamp = time.Now().Local().Format("2006-01-02 15:04:05")
			cdr.EndStamp = time.Now().Local().Format("2006-01-02 15:04:05")
			cdr.CallerIDName = channel.Data.CidName
			cdr.CallerIDNumber = channel.Data.CidNumber
			cdr.DestinationNumber = channel.Data.DestNumber
			cdr.Session = channel.Data.Params["sip_h_X-FS-Session"]

			params := &RequestParam{
				NodeUUID: boy.NodeUUID,
				NodeIP:   "127.0.0.1",
				UUID:     channel.Data.Uuid,
				CDR:      cdr,
			}

			rpc := *&ctrl.Request{
				Version: "2.0",
				Method:  "Event.CDR",
				Params:  ctrl.ToRawMessage(params),
			}

			str, _ := json.MarshalIndent(rpc, "", "  ")
			ctrl.Publish("cn.xswitch.event.cdr", str)
		}
		return
	}

}

func (boy *TBoyACD) Dial(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.DialRequest

	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		return
	}

	controller := "cn.xswitch.ctrl." + request.CtrlUuid

	uuid := request.Destination.CallParams[0].Uuid

	call_params_params := request.Destination.CallParams[0].Params

	if call_params_params == nil {
		call_params_params = make(map[string]string)
	}

	channelEvent := xctrl.ChannelEvent{
		NodeUuid:    boy.NodeUUID,
		Uuid:        uuid,
		Direction:   "outbound",
		State:       "CALLING",
		CidName:     "TEST",
		CidNumber:   request.Destination.CallParams[0].CidNumber,
		DestNumber:  request.Destination.CallParams[0].DestNumber,
		AnswerEpoch: uint32(time.Now().Unix()),
		Answered:    true,
		CreateEpoch: uint32(time.Now().Unix()),
		Params:      call_params_params,
	}

	channelEvent.Params["xcc_session"] = channelEvent.Params["sip_h_X-FS-Session"]
	channel := &FakeChannel{
		CtrlUuid: request.CtrlUuid,
		Data:     &channelEvent,
	}
	CacheChannel(uuid, channel)

	//主叫
	if call_params_params["xcc_action"] != "" {
		event_req := ctrl.Request{
			Version: "2.0",
			Method:  "Event.Channel",
			Params:  ctrl.ToRawMessage(channelEvent),
		}

		req_str, _ := json.MarshalIndent(event_req, "", "  ")
		ctrl.Publish(controller, req_str)

		channelEvent.State = "ANSWERED"
		event_req.Params = ctrl.ToRawMessage(channelEvent)
		req_str, _ = json.MarshalIndent(event_req, "", "  ")
		ctrl.Publish(controller, req_str)

		channelEvent.State = "READY"
		event_req.Params = ctrl.ToRawMessage(channelEvent)
		req_str, _ = json.MarshalIndent(event_req, "", "  ")
		ctrl.Publish(controller, req_str)

		channel = &FakeChannel{
			CtrlUuid: request.CtrlUuid,
			Data:     &channelEvent,
		}
		if channel.Context == nil {
			channel.Context, channel.Cancel = context.WithCancel(context.Background())
		}

		CacheChannel(uuid, channel)

		if channelEvent.GetParams()["xcc_action"] == "" {
			go func() {
				t := time.NewTimer(3 * time.Second)
				select {
				case <-channel.Context.Done():
					return
				case <-t.C:
					controller := "cn.xswitch.ctrl." + channel.CtrlUuid

					channelEvent.State = "DESTROY"
					channelEvent.Cause = "NORMAL_CLEARING"
					event_req.Params = ctrl.ToRawMessage(channelEvent)
					req_str, _ = json.MarshalIndent(event_req, "", "  ")
					ctrl.Publish(controller, req_str)

					cdr := new(CallDetailRecord)
					err = json.Unmarshal([]byte(cdr_template), cdr)

					if err != nil {
						log.Error(err)
					}

					if send_cdr {

						cdr.UUID = channel.Data.Uuid
						cdr.Leg = "a"
						cdr.StartStamp = time.Now().Local().Format("2006-01-02 15:04:05")
						cdr.EndStamp = time.Now().Local().Format("2006-01-02 15:04:05")
						cdr.CallerIDName = channel.Data.CidName
						cdr.CallerIDNumber = channel.Data.CidNumber
						cdr.DestinationNumber = channel.Data.DestNumber
						cdr.Session = channel.Data.Params["sip_h_X-FS-Session"]

						params := &RequestParam{
							NodeUUID: "test-node",
							NodeIP:   "127.0.0.1",
							UUID:     channel.Data.Uuid,
							CDR:      cdr,
						}

						rpc := *&ctrl.Request{
							Version: "2.0",
							Method:  "Event.CDR",
							Params:  ctrl.ToRawMessage(params),
						}

						str, _ := json.MarshalIndent(rpc, "", "  ")
						ctrl.Publish("cn.xswitch.event.cdr", str)
					}
				}

			}()
		}

		// response
		res := xctrl.Response{
			Code:     200,
			Message:  "OK",
			NodeUuid: boy.NodeUUID,
		}

		res_bytes, _ := json.MarshalIndent(res, "", "  ")
		raw := json.RawMessage(res_bytes)

		result := ctrl.Result{
			Version: "2.0",
			ID:      msg.ID,
			Result:  &raw,
		}

		result_bytes, _ := json.MarshalIndent(result, "", "  ")

		if reply == "" {
			reply = "cn.xswitch.ctrl." + request.CtrlUuid
		}

		ctrl.Publish(reply, result_bytes)
		return
	}
	// 被叫不接
	if !boy.Options.PeerAnswer {
		event_req := ctrl.Request{
			Version: "2.0",
			Method:  "Event.Channel",
			Params:  ctrl.ToRawMessage(channelEvent),
		}

		req_str, _ := json.MarshalIndent(event_req, "", "  ")
		ctrl.Publish(controller, req_str)
		time.Sleep(time.Second * 5)

		channelEvent.State = "DESTROY"
		channelEvent.Cause = "NORMAL_CLEARING"
		event_req.Params = ctrl.ToRawMessage(channelEvent)
		req_str, _ = json.MarshalIndent(event_req, "", "  ")
		ctrl.Publish(controller, req_str)

	} else {
		event_req := ctrl.Request{
			Version: "2.0",
			Method:  "Event.Channel",
			Params:  ctrl.ToRawMessage(channelEvent),
		}

		req_str, _ := json.MarshalIndent(event_req, "", "  ")
		ctrl.Publish(controller, req_str)

		channelEvent.State = "ANSWERED"
		event_req.Params = ctrl.ToRawMessage(channelEvent)
		req_str, _ = json.MarshalIndent(event_req, "", "  ")
		ctrl.Publish(controller, req_str)

		channelEvent.State = "READY"
		event_req.Params = ctrl.ToRawMessage(channelEvent)
		req_str, _ = json.MarshalIndent(event_req, "", "  ")
		ctrl.Publish(controller, req_str)
	}

}
