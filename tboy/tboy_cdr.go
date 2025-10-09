package tboy

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"git.xswitch.cn/xswitch/proto/go/proto/xctrl"
	"git.xswitch.cn/xswitch/xctrl/ctrl"
	"git.xswitch.cn/xswitch/xctrl/ctrl/nats"
)

type FakeFakeCdrChannel struct {
	CtrlUuid string
	UUID     string
	Lock     sync.RWMutex
	Data     *xctrl.ChannelEvent
	Context  context.Context
	Cancel   context.CancelFunc
	// store cdr data
	CDRData *CallDetailRecord
}

type FakeFakeCdrServer struct {
	NodeUUID string
	Domain   string
	Channels map[string]*FakeFakeCdrChannel
	Lock     sync.RWMutex
}

func NewFakeFakeCdrServer(nodeUUID, domain string) *FakeFakeCdrServer {
	return &FakeFakeCdrServer{
		NodeUUID: nodeUUID,
		Domain:   domain,
		Channels: make(map[string]*FakeFakeCdrChannel),
	}
}

func (server *FakeFakeCdrServer) CacheChannel(uuid string, channel *FakeFakeCdrChannel) {
	server.Lock.Lock()
	defer server.Lock.Unlock()
	server.Channels[uuid] = channel
}

func (server *FakeFakeCdrServer) GetChannel(uuid string) (*FakeFakeCdrChannel, bool) {
	server.Lock.RLock()
	defer server.Lock.RUnlock()
	channel, ok := server.Channels[uuid]
	return channel, ok
}

func (server *FakeFakeCdrServer) DeleteChannel(uuid string) {
	server.Lock.Lock()
	defer server.Lock.Unlock()
	delete(server.Channels, uuid)
}

func (server *FakeFakeCdrServer) OK(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.Request

	err := json.Unmarshal(*msg.Params, &request)
	if err != nil {
		return
	}

	res := xctrl.Response{
		Code:     200,
		Message:  "OK",
		NodeUuid: server.NodeUUID,
		Uuid:     request.Uuid,
	}

	resBytes, _ := json.Marshal(res)
	raw := json.RawMessage(resBytes)

	result := ctrl.Result{
		Version: "2.0",
		ID:      msg.ID,
		Result:  &raw,
	}

	resultBytes, _ := json.Marshal(result)

	if reply == "" {
		reply = "cn.xswitch.ctrl." + request.CtrlUuid
	}

	ctrl.Publish(reply, resultBytes)
}

func (server *FakeFakeCdrServer) Error(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.Request

	err := json.Unmarshal(*msg.Params, &request)
	if err != nil {
		return
	}

	res := xctrl.Response{
		Code:     404,
		Message:  "Unsupported Method " + msg.Method,
		NodeUuid: server.NodeUUID,
		Uuid:     request.Uuid,
	}

	resBytes, _ := json.Marshal(res)
	raw := json.RawMessage(resBytes)

	result := ctrl.Result{
		Version: "2.0",
		ID:      msg.ID,
		Result:  &raw,
	}

	resultBytes, _ := json.Marshal(result)

	if reply == "" {
		reply = "cn.xswitch.ctrl." + request.CtrlUuid
	}

	ctrl.Publish(reply, resultBytes)
}

func (server *FakeFakeCdrServer) Answer(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.Request
	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		log.Error(err)
		return
	}

	channel, ok := server.GetChannel(request.Uuid)
	if ok {
		channel.CtrlUuid = request.CtrlUuid
	}

	server.OK(ctx, msg, reply)

	if channel != nil {
		controller := "cn.xswitch.ctrl." + channel.CtrlUuid
		channel.Data.State = "ANSWERED"
		channel.Data.AnswerEpoch = uint32(time.Now().Unix())

		eventReq := ctrl.Request{
			Version: "2.0",
			Method:  "Event.Channel",
			Params:  ctrl.ToRawMessage(channel.Data),
		}

		reqStr, _ := json.Marshal(eventReq)
		ctrl.Publish(controller, reqStr)

		if channel.CDRData != nil {
			channel.CDRData.AnswerStamp = time.Now().Format("2006-01-02 15:04:05")
		}
	}
}

func (server *FakeFakeCdrServer) Hangup(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.Request
	err := json.Unmarshal(*msg.Params, &request)

	if err != nil {
		log.Error(err)
		return
	}

	channel, ok := server.GetChannel(request.Uuid)
	if !ok {
		server.Error(ctx, msg, reply)
		return
	}

	if channel.Cancel != nil {
		channel.Cancel()
	}

	server.OK(ctx, msg, reply)

	// send DESTROY event
	controller := "cn.xswitch.ctrl." + channel.CtrlUuid
	channel.Data.State = "DESTROY"
	channel.Data.Cause = "NORMAL_CLEARING"

	eventReq := ctrl.Request{
		Version: "2.0",
		Method:  "Event.Channel",
		Params:  ctrl.ToRawMessage(channel.Data),
	}

	reqStr, _ := json.Marshal(eventReq)
	ctrl.Publish(controller, reqStr)

	// send CDR
	server.SendCDR(channel)
	server.DeleteChannel(request.Uuid)
}

func (server *FakeFakeCdrServer) SendCDR(channel *FakeFakeCdrChannel) {
	if channel.CDRData == nil {
		return
	}

	channel.CDRData.EndStamp = time.Now().Format("2006-01-02 15:04:05")

	startTime, _ := time.Parse("2006-01-02 15:04:05", channel.CDRData.StartStamp)
	endTime, _ := time.Parse("2006-01-02 15:04:05", channel.CDRData.EndStamp)
	duration := int(endTime.Sub(startTime).Seconds())
	channel.CDRData.Duration = duration

	if channel.CDRData.AnswerStamp != "" {
		answerTime, _ := time.Parse("2006-01-02 15:04:05", channel.CDRData.AnswerStamp)
		billsec := int(endTime.Sub(answerTime).Seconds())
		channel.CDRData.Billsec = billsec
	}

	params := &RequestParam{
		NodeUUID: server.NodeUUID,
		NodeIP:   "127.0.0.1",
		UUID:     channel.Data.Uuid,
		CDR:      channel.CDRData,
	}

	rpc := ctrl.Request{
		Version: "2.0",
		Method:  "Event.CDR",
		Params:  ctrl.ToRawMessage(params),
	}

	str, _ := json.Marshal(rpc)
	ctrl.Publish("cn.xswitch.cdr", str)
	log.Infof("Sent CDR for channel: %s", channel.Data.Uuid)
}

func (server *FakeFakeCdrServer) Dial(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.DialRequest

	err := json.Unmarshal(*msg.Params, &request)
	if err != nil {
		log.Error(err)
		return
	}

	controller := "cn.xswitch.ctrl." + request.CtrlUuid

	if len(request.Destination.CallParams) < 1 {
		server.Error(ctx, msg, reply)
		return
	}

	uuid := request.Destination.CallParams[0].Uuid
	callParams := request.Destination.CallParams[0].Params

	if callParams == nil {
		callParams = make(map[string]string)
	}

	channelEvent := xctrl.ChannelEvent{
		NodeUuid:    server.NodeUUID,
		Uuid:        uuid,
		Direction:   "outbound",
		State:       "CALLING",
		CidName:     "TEST",
		CidNumber:   request.Destination.CallParams[0].CidNumber,
		DestNumber:  request.Destination.CallParams[0].DestNumber,
		AnswerEpoch: uint32(time.Now().Unix()),
		Answered:    false,
		CreateEpoch: uint32(time.Now().Unix()),
		Params:      callParams,
	}

	channelEvent.Params["xcc_session"] = channelEvent.Params["sip_h_X-FS-Session"]

	cdr := &CallDetailRecord{
		UUID:              uuid,
		Domain:            server.Domain,
		Context:           "default",
		CallerIDName:      channelEvent.CidName,
		CallerIDNumber:    channelEvent.CidNumber,
		DestinationNumber: channelEvent.DestNumber,
		StartStamp:        time.Now().Format("2006-01-02 15:04:05"),
		HangupCause:       "NORMAL_CLEARING",
		Leg:               "a",
		Direction:         "outbound",
		Session:           channelEvent.Params["sip_h_X-FS-Session"],
	}

	channel := &FakeFakeCdrChannel{
		CtrlUuid: request.CtrlUuid,
		UUID:     uuid,
		Data:     &channelEvent,
		CDRData:  cdr,
	}
	channel.Context, channel.Cancel = context.WithCancel(context.Background())

	server.CacheChannel(uuid, channel)

	eventReq := ctrl.Request{
		Version: "2.0",
		Method:  "Event.Channel",
		Params:  ctrl.ToRawMessage(channelEvent),
	}

	reqStr, _ := json.Marshal(eventReq)
	ctrl.Publish(controller, reqStr)

	go server.simulateCall(channel)

	server.OK(ctx, msg, reply)
}

func (server *FakeFakeCdrServer) simulateCall(channel *FakeFakeCdrChannel) {
	controller := "cn.xswitch.ctrl." + channel.CtrlUuid

	time.Sleep(1 * time.Second)

	channel.Data.State = "ANSWERED"
	channel.Data.AnswerEpoch = uint32(time.Now().Unix())
	channel.Data.Answered = true

	eventReq := ctrl.Request{
		Version: "2.0",
		Method:  "Event.Channel",
		Params:  ctrl.ToRawMessage(channel.Data),
	}

	reqStr, _ := json.Marshal(eventReq)
	ctrl.Publish(controller, reqStr)

	if channel.CDRData != nil {
		channel.CDRData.AnswerStamp = time.Now().Format("2006-01-02 15:04:05")
	}

	// wait 3s and send hangup event
	time.Sleep(3 * time.Second)

	select {
	case <-channel.Context.Done():
		return
	default:
		channel.Data.State = "DESTROY"
		channel.Data.Cause = "NORMAL_CLEARING"

		eventReq.Params = ctrl.ToRawMessage(channel.Data)
		reqStr, _ = json.Marshal(eventReq)
		ctrl.Publish(controller, reqStr)

		server.SendCDR(channel)

		server.DeleteChannel(channel.UUID)
	}
}

func (server *FakeFakeCdrServer) Event(msg *ctrl.Message, natsEvent nats.Event) {
	topic := natsEvent.Topic()
	reply := natsEvent.Reply()
	log.Infof("Received event: %s %s", topic, msg.Method)

	if msg.Method == "" && msg.Result != nil {
		log.Infof("Got a response: %s", msg.ID)
		return
	}

	ctx := context.Background()

	switch msg.Method {
	case "XNode.Answer":
		server.Answer(ctx, msg, reply)
	case "XNode.Hangup":
		server.Hangup(ctx, msg, reply)
	case "XNode.Dial":
		server.Dial(ctx, msg, reply)
	default:
		log.Errorf("Unsupported Method: %s", msg.Method)
		server.Error(ctx, msg, reply)
	}
}
