package tboy

import (
	"context"
	"encoding/json"
	"sync"

	"git.xswitch.cn/xswitch/proto/go/proto/xctrl"
	"git.xswitch.cn/xswitch/xctrl/ctrl"
	"git.xswitch.cn/xswitch/xctrl/ctrl/nats"
)

type FakeLayoutChannel struct {
	*FakeChannel
	ConferenceName string
	Layout         string
	AgoraUID       string
}

type Layout struct {
	*TBoy
	LayoutChannels map[string]*FakeLayoutChannel
	LayoutLock     sync.RWMutex
}

func NewTBoyLayout(tboy *TBoy) *Layout {
	return &Layout{
		TBoy:           tboy,
		LayoutChannels: make(map[string]*FakeLayoutChannel),
	}
}

type LayoutEvent struct {
	NodeUuid       string `json:"node_uuid"`
	ConferenceName string `json:"conference_name"`
	AgoraUID       string `json:"agora_uid,omitempty"`
	Action         string `json:"action"`
	Layout         string `json:"layout,omitempty"`
	Value          string `json:"value,omitempty"`
}

type ConfControlRequest struct {
	Command        string `json:"command"`
	ConferenceName string `json:"conference_name"`
	AgoraUID       string `json:"agora_uid,omitempty"`
	Value          string `json:"value,omitempty"`
}

type ConfControlResponse struct {
	Code     int32  `json:"code"`
	Message  string `json:"message"`
	NodeUuid string `json:"node_uuid"`
}

func (boy *Layout) Init() {
	boy.TBoy.Init()
	boy.LayoutChannels = make(map[string]*FakeLayoutChannel)
}

func (boy *Layout) CacheLayoutChannel(uuid string, channel *FakeLayoutChannel) {
	boy.LayoutLock.Lock()
	defer boy.LayoutLock.Unlock()
	boy.LayoutChannels[uuid] = channel
}

func (boy *Layout) GetLayoutChannel(uuid string) (*FakeLayoutChannel, bool) {
	boy.LayoutLock.RLock()
	defer boy.LayoutLock.RUnlock()
	channel, ok := boy.LayoutChannels[uuid]
	return channel, ok
}

func (boy *Layout) ConfControl(ctx context.Context, msg *ctrl.Message, reply string) {
	var req ConfControlRequest
	err := json.Unmarshal(*msg.Params, &req)
	if err != nil {
		log.Error("Failed to unmarshal confControl request: ", err)
		boy.ConfControlError(ctx, msg, reply, "Invalid request parameters")
		return
	}

	var layoutChannel *FakeLayoutChannel
	for _, channel := range boy.LayoutChannels {
		if channel.ConferenceName == req.ConferenceName {
			if req.AgoraUID == "" || channel.AgoraUID == req.AgoraUID {
				layoutChannel = channel
				break
			}
		}
	}

	if layoutChannel == nil {
		log.Errorf("Conference not found: %s, AgoraUID: %s", req.ConferenceName, req.AgoraUID)
		boy.ConfControlError(ctx, msg, reply, "Conference not found")
		return
	}

	switch req.Command {
	case "hangup":
		boy.handleHangup(layoutChannel, req)
	case "vid-res-id":
		boy.handleVidResId(layoutChannel, req)
	case "layout":
		boy.handleLayout(layoutChannel, req)
	case "set-layout":
		boy.handleSetLayout(layoutChannel, req)
	default:
		log.Errorf("Unsupported confControl command: %s", req.Command)
		boy.ConfControlError(ctx, msg, reply, "Unsupported command")
		return
	}

	boy.sendLayoutEvent(layoutChannel, req.Command, req.Value)

	boy.ConfControlOK(ctx, msg, reply)
}

func (boy *Layout) handleHangup(channel *FakeLayoutChannel, req ConfControlRequest) {
	log.Infof("Hangup conference: %s, AgoraUID: %s", req.ConferenceName, req.AgoraUID)

	boy.RemoveFromLayout(channel.FakeChannel.Data.Uuid)

	boy.sendLayoutEvent(channel, "hangup", "")
}

func (boy *Layout) handleVidResId(channel *FakeLayoutChannel, req ConfControlRequest) {
	log.Infof("Video resolution ID for conference: %s, value: %s", req.ConferenceName, req.Value)

	if req.Value == "all clear" {
		channel.Layout = "default"
	} else {
		channel.Layout = req.Value
	}
}

func (boy *Layout) handleLayout(channel *FakeLayoutChannel, req ConfControlRequest) {
	log.Infof("Get layout for conference: %s", req.ConferenceName)

	if req.Value != "" {
		channel.Layout = req.Value
	}
}

func (boy *Layout) handleSetLayout(channel *FakeLayoutChannel, req ConfControlRequest) {
	log.Infof("Set layout for conference: %s to %s", req.ConferenceName, req.Value)

	if req.Value != "" {
		channel.Layout = req.Value
	}
}

func (boy *Layout) sendLayoutEvent(channel *FakeLayoutChannel, action, value string) {
	event := LayoutEvent{
		NodeUuid:       boy.NodeUUID,
		ConferenceName: channel.ConferenceName,
		AgoraUID:       channel.AgoraUID,
		Action:         action,
		Layout:         channel.Layout,
		Value:          value,
	}

	eventReq := ctrl.Request{
		Version: "2.0",
		Method:  "Event.Layout",
		Params:  ctrl.ToRawMessage(event),
	}

	eventBytes, _ := json.Marshal(eventReq)
	ctrl.Publish("cn.xswitch.event.layout", eventBytes)
}

func (boy *Layout) ConfControlOK(ctx context.Context, msg *ctrl.Message, reply string) {
	res := ConfControlResponse{
		Code:     200,
		Message:  "OK",
		NodeUuid: boy.NodeUUID,
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
		reply = "cn.xswitch.cman.control"
	}

	ctrl.Publish(reply, resultBytes)
}

func (boy *Layout) ConfControlError(ctx context.Context, msg *ctrl.Message, reply string, errorMsg string) {
	res := ConfControlResponse{
		Code:     500,
		Message:  errorMsg,
		NodeUuid: boy.NodeUUID,
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
		reply = "cn.xswitch.cman.control"
	}

	ctrl.Publish(reply, resultBytes)
}

func (boy *Layout) JoinConference(uuid, conferenceName, agoraUID string) {
	baseChannel, ok := Channels[uuid]
	if !ok {
		log.Error("Channel not found for conference join: ", uuid)
		return
	}

	layoutChannel := &FakeLayoutChannel{
		FakeChannel:    baseChannel,
		ConferenceName: conferenceName,
		AgoraUID:       agoraUID,
		Layout:         "default", // 默认布局
	}

	boy.CacheLayoutChannel(uuid, layoutChannel)

	boy.sendLayoutEvent(layoutChannel, "join", "")
}

func (boy *Layout) RemoveFromLayout(uuid string) {
	boy.LayoutLock.Lock()
	defer boy.LayoutLock.Unlock()

	if channel, ok := boy.LayoutChannels[uuid]; ok {
		boy.sendLayoutEvent(channel, "leave", channel.ConferenceName)
		delete(boy.LayoutChannels, uuid)
	}
}

func (boy *Layout) GetConferenceLayout(conferenceName string) string {
	boy.LayoutLock.RLock()
	defer boy.LayoutLock.RUnlock()

	for _, channel := range boy.LayoutChannels {
		if channel.ConferenceName == conferenceName {
			return channel.Layout
		}
	}
	return "default"
}

func (boy *Layout) SetConferenceLayout(conferenceName, layout string) {
	boy.LayoutLock.Lock()
	defer boy.LayoutLock.Unlock()

	for _, channel := range boy.LayoutChannels {
		if channel.ConferenceName == conferenceName {
			channel.Layout = layout
			boy.sendLayoutEvent(channel, "layout-change", layout)
		}
	}
}

func (boy *Layout) Event(msg *ctrl.Message, natsEvent nats.Event) {
	topic := natsEvent.Topic()
	reply := natsEvent.Reply()
	log.Infof("%s %s", topic, msg.Method)

	if msg.Method == "" && msg.Result != nil {
		log.Infof("Got a response: %s", msg.ID)
		return
	}

	ctx := context.Background()

	switch msg.Method {
	case "confControl":
		boy.ConfControl(ctx, msg, reply)
	default:
		boy.TBoy.Event(msg, natsEvent)
	}
}

func (boy *Layout) Hangup(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.Request
	err := json.Unmarshal(*msg.Params, &request)
	if err != nil {
		return
	}

	boy.RemoveFromLayout(request.Uuid)

	boy.TBoy.Hangup(ctx, msg, reply)
}

func (boy *Layout) GetConferenceMembers(conferenceName string) []*FakeLayoutChannel {
	boy.LayoutLock.RLock()
	defer boy.LayoutLock.RUnlock()

	var members []*FakeLayoutChannel
	for _, channel := range boy.LayoutChannels {
		if channel.ConferenceName == conferenceName {
			members = append(members, channel)
		}
	}
	return members
}

func (boy *Layout) GetConferenceStats(conferenceName string) map[string]interface{} {
	members := boy.GetConferenceMembers(conferenceName)

	stats := map[string]interface{}{
		"conference_name": conferenceName,
		"member_count":    len(members),
		"current_layout":  boy.GetConferenceLayout(conferenceName),
		"agora_uids":      []string{},
	}

	agoraUIDs := make([]string, 0, len(members))
	for _, member := range members {
		if member.AgoraUID != "" {
			agoraUIDs = append(agoraUIDs, member.AgoraUID)
		}
	}
	stats["agora_uids"] = agoraUIDs

	return stats
}
