package tboy

import (
	"context"
	"encoding/json"
	"git.xswitch.cn/xswitch/xctrl/ctrl/nats"
	"sync"

	"git.xswitch.cn/xswitch/proto/go/proto/xctrl"
	"git.xswitch.cn/xswitch/xctrl/ctrl"
	"github.com/sirupsen/logrus"
)

func init() {
	log = logrus.New()
	log.SetReportCaller(true)
}

type FakeConferenceChannel struct {
	*FakeChannel
	ConferenceId string
	Muted        bool
	VMuted       bool
}

type TBoyConference struct {
	*TBoy
	ConferenceChannels map[string]*FakeConferenceChannel
	ConfLock           sync.RWMutex
}

func NewTBoyConference(tboy *TBoy) *TBoyConference {
	return &TBoyConference{
		TBoy:               tboy,
		ConferenceChannels: make(map[string]*FakeConferenceChannel),
	}
}

type ConferenceEvent struct {
	NodeUuid   string `json:"node_uuid"`
	Uuid       string `json:"uuid"`
	Action     string `json:"action"`
	Conference string `json:"conference"`
	Muted      bool   `json:"muted,omitempty"`
	VMuted     bool   `json:"vmute,omitempty"`
}

func (boy *TBoyConference) Init() {
	boy.TBoy.Init()
	boy.ConferenceChannels = make(map[string]*FakeConferenceChannel)
}

func (boy *TBoyConference) CacheConferenceChannel(uuid string, channel *FakeConferenceChannel) {
	boy.ConfLock.Lock()
	defer boy.ConfLock.Unlock()
	boy.ConferenceChannels[uuid] = channel
}

func (boy *TBoyConference) GetConferenceChannel(uuid string) (*FakeConferenceChannel, bool) {
	boy.ConfLock.RLock()
	defer boy.ConfLock.RUnlock()
	channel, ok := boy.ConferenceChannels[uuid]
	return channel, ok
}

func (boy *TBoyConference) Conference(ctx context.Context, msg *ctrl.Message, reply string) {
	var req xctrl.ConferenceRequest
	err := json.Unmarshal(*msg.Params, &req)
	if err != nil {
		log.Error("Failed to unmarshal conference request: ", err)
		boy.Error(ctx, msg, reply)
		return
	}

	baseChannel, ok := Channels[req.Uuid]
	if !ok {
		log.Error("Channel not found: ", req.Uuid)
		boy.Error(ctx, msg, reply)
		return
	}

	channel, ok := boy.GetConferenceChannel(req.Uuid)
	if !ok {
		channel = &FakeConferenceChannel{
			FakeChannel:  baseChannel,
			ConferenceId: "",
			Muted:        false,
			VMuted:       false,
		}
		boy.CacheConferenceChannel(req.Uuid, channel)
	}

	action := "join"
	var muted, VMuted bool

	for _, flag := range req.Flags {
		switch flag {
		case "mute":
			action = "mute"
			muted = true
		case "unmute":
			action = "unmute"
			muted = false
		case "vmute":
			action = "vmute"
			VMuted = true
		case "unvmute":
			action = "unvmute"
			VMuted = false
		case "deaf":
			action = "deaf"
			muted = true
		case "undeaf":
			action = "undeaf"
			muted = false
		}
	}

	if action == "join" {
		channel.ConferenceId = req.Name
		channel.Muted = false
		channel.VMuted = false
	}

	if action == "mute" || action == "unmute" || action == "deaf" || action == "undeaf" {
		channel.Muted = muted
	}

	if action == "vmute" || action == "unvmute" {
		channel.VMuted = VMuted
	}

	// 发送会议事件
	boy.sendConferenceEvent(channel, action, req.Name, channel.Muted, channel.VMuted)
	boy.OK(ctx, msg, reply)
}

func (boy *TBoyConference) sendConferenceEvent(channel *FakeConferenceChannel, action, confId string, muted, VMuted bool) {
	event := ConferenceEvent{
		NodeUuid:   boy.NodeUUID,
		Uuid:       channel.Data.Uuid,
		Action:     action,
		Conference: confId,
		Muted:      muted,
		VMuted:     VMuted,
	}

	eventReq := ctrl.Request{
		Version: "2.0",
		Method:  "Event.Conference",
		Params:  ctrl.ToRawMessage(event),
	}

	eventBytes, _ := json.Marshal(eventReq)
	ctrl.Publish("cn.xswitch.event.conference", eventBytes)
}

func (boy *TBoyConference) ConferenceInfo(ctx context.Context, msg *ctrl.Message, reply string) {
	var req xctrl.ConferenceRequest
	err := json.Unmarshal(*msg.Params, &req)
	if err != nil {
		log.Error("Failed to unmarshal conference info request: ", err)
		boy.Error(ctx, msg, reply)
		return
	}

	conferenceInfo := &xctrl.ConferenceInfo{
		ConferenceName: req.Name,
		ConferenceUuid: "conf-" + req.Name + "-uuid",
		MemberCount:    boy.getConferenceMemberCount(req.Name),
		Running:        true,
		Answered:       true,
		Domain:         boy.Domain,
	}

	res := xctrl.ConferenceInfoResponse{
		Code:     200,
		Message:  "OK",
		NodeUuid: boy.NodeUUID,
		Data:     conferenceInfo,
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
		if channel, ok := Channels[req.Uuid]; ok {
			reply = "cn.xswitch.ctrl." + channel.CtrlUuid
		} else {
			reply = "cn.xswitch.ctrl.default"
		}
	}

	ctrl.Publish(reply, resultBytes)
}

func (boy *TBoyConference) getConferenceMemberCount(conferenceName string) int32 {
	boy.ConfLock.RLock()
	defer boy.ConfLock.RUnlock()

	count := 0
	for _, channel := range boy.ConferenceChannels {
		if channel.ConferenceId == conferenceName {
			count++
		}
	}
	return int32(count)
}

func (boy *TBoyConference) GetConferenceMembers(conferenceName string) []*FakeConferenceChannel {
	boy.ConfLock.RLock()
	defer boy.ConfLock.RUnlock()

	var members []*FakeConferenceChannel
	for _, channel := range boy.ConferenceChannels {
		if channel.ConferenceId == conferenceName {
			members = append(members, channel)
		}
	}
	return members
}

func (boy *TBoyConference) RemoveFromConference(uuid string) {
	boy.ConfLock.Lock()
	defer boy.ConfLock.Unlock()

	if channel, ok := boy.ConferenceChannels[uuid]; ok {
		// leave conf
		boy.sendConferenceEvent(channel, "leave", channel.ConferenceId, channel.Muted, channel.VMuted)
		delete(boy.ConferenceChannels, uuid)
	}
}

func (boy *TBoyConference) Event(msg *ctrl.Message, natsEvent nats.Event) {
	topic := natsEvent.Topic()
	reply := natsEvent.Reply()
	log.Infof("%s %s", topic, msg.Method)

	if msg.Method == "" && msg.Result != nil {
		log.Infof("Got a response: %s", msg.ID)
		return
	}

	ctx := context.Background()

	switch msg.Method {
	case "XNode.Conference":
		boy.Conference(ctx, msg, reply)
	case "XNode.ConferenceInfo":
		boy.ConferenceInfo(ctx, msg, reply)
	default:
		// call father func
		boy.TBoy.Event(msg, natsEvent)
	}
}

func (boy *TBoyConference) Hangup(ctx context.Context, msg *ctrl.Message, reply string) {
	var request xctrl.Request
	err := json.Unmarshal(*msg.Params, &request)
	if err != nil {
		return
	}

	// rm from conf
	boy.RemoveFromConference(request.Uuid)

	boy.TBoy.Hangup(ctx, msg, reply)
}
