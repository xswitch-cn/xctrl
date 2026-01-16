package ctrl

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/xswitch-cn/proto/go/proto/xctrl"
	"github.com/xswitch-cn/proto/xctrl/util/log"
	"github.com/xswitch-cn/xctrl/ctrl/nats"
)

const (
	nodeUUID = "test.node-uuid"
)

func init() {
	natsURL := os.Getenv("NATS_ADDRESS")

	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	err := Init(true, natsURL)
	if err != nil {
		log.Error(err)
	}
}

func TestPlayWithTimeout1(t *testing.T) {
	channelEvent := &xctrl.ChannelEvent{
		NodeUuid: nodeUUID,
	}

	channel := &Channel{
		ChannelEvent: channelEvent,
		CtrlUuid:     UUID(),
	}

	channel.NodeUuid = nodeUUID

	req := &xctrl.PlayRequest{
		CtrlUuid: UUID(),
		Uuid:     "test-uuid",
		Media: &xctrl.Media{
			Data: "/tmp/test.wav",
		},
	}

	sub, err := Subscribe("cn.xswitch.node."+nodeUUID, func(c context.Context, e nats.Event) error {
		return nil
	}, nodeUUID)
	if err != nil {
		t.Error(err)
	}

	res := channel.PlayWithTimeout(req, 100*time.Millisecond)

	fmt.Println(res.Code)
	if res.Code != 408 {
		t.Error(res)
	}

	sub.Unsubscribe()
}

func TestPlayWithTimeout2(t *testing.T) {
	channelEvent := &xctrl.ChannelEvent{
		NodeUuid: nodeUUID,
	}

	channel := &Channel{
		ChannelEvent: channelEvent,
		CtrlUuid:     UUID(),
	}

	channel.NodeUuid = nodeUUID

	req := &xctrl.PlayRequest{
		CtrlUuid: UUID(),
		Uuid:     "test-uuid",
		Media: &xctrl.Media{
			Data: "/tmp/test.wav",
		},
	}

	sub, err := Subscribe("cn.xswitch.node."+nodeUUID, func(c context.Context, e nats.Event) error {
		var request Request
		err := json.Unmarshal(e.Message().Body, &request)
		if err != nil {
			t.Error(err)
		}

		response := &xctrl.Response{
			Code:     200,
			Message:  "OK",
			NodeUuid: nodeUUID,
			Uuid:     "test-uuid",
		}

		rpc := &Response{
			Version: "2.0",
			ID:      request.ID,
			Result:  ToRawMessage(response),
		}

		err = PublishJSON(e.Reply(), rpc)
		if err != nil {
			t.Error(err)
		}
		return nil
	}, nodeUUID)
	if err != nil {
		return
	}

	res := channel.PlayWithTimeout(req, 100*time.Millisecond)

	if res.Code != 200 {
		t.Error(res)
	}

	res = channel.Play(req)

	if res.Code != 200 {
		t.Error(res)
	}

	sub.Unsubscribe()
}

func TestFIFO(t *testing.T) {
	// subject := "cn.xswitch.ctrl"
	//获取nats地址
	natsURL := os.Getenv("NATS_ADDRESS")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	//初始化 ctrl
	err := Init(true, natsURL)
	if err != nil {
		t.Error(err)
	}

	nodeUUID := "test.node-uuid"

	//订阅主题
	_, err = Subscribe("cn.xswitch.node."+nodeUUID, func(c context.Context, e nats.Event) error {
		var request Request
		err := json.Unmarshal(e.Message().Body, &request)
		if err != nil {
			t.Error(err)
		}

		response := &xctrl.Response{
			Code:     200,
			Message:  "OK",
			NodeUuid: nodeUUID,
			Uuid:     "test-uuid",
		}

		rpc := &Response{
			Version: "2.0",
			ID:      request.ID,
			Result:  ToRawMessage(response),
		}
		err = PublishJSON(e.Reply(), rpc)
		if err != nil {
			t.Error(err)
		}
		return nil
	}, nodeUUID)
	if err != nil {
		t.Error(err)
	}

	channelEvent := &xctrl.ChannelEvent{
		NodeUuid: nodeUUID,
	}

	channel := &Channel{
		ChannelEvent: channelEvent,
		CtrlUuid:     UUID(),
	}

	req := &xctrl.FIFORequest{
		Uuid:         UUID(),
		Name:         "test_name",
		Inout:        "out",
		WaitMusic:    "/tmp/test.wav",
		ExitAnnounce: "/tmp/test.wav",
		Priority:     0,
	}

	res := channel.FIFO(req)

	if res.Code != 200 {
		t.Error(res)
	}

}

func TestChannel_Callcenter(t *testing.T) {
	// subject := "cn.xswitch.ctrl"
	//获取nats地址
	natsURL := os.Getenv("NATS_ADDRESS")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	//初始化 ctrl
	err := Init(true, natsURL)
	if err != nil {
		t.Error(err)
	}

	nodeUUID := "test.node-uuid"

	//订阅主题
	_, err = Subscribe("cn.xswitch.node."+nodeUUID, func(c context.Context, e nats.Event) error {
		var request Request
		err := json.Unmarshal(e.Message().Body, &request)
		if err != nil {
			t.Error(err)
		}

		response := &xctrl.Response{
			Code:     200,
			Message:  "OK",
			NodeUuid: nodeUUID,
			Uuid:     "test-uuid",
		}

		rpc := &Response{
			Version: "2.0",
			ID:      request.ID,
			Result:  ToRawMessage(response),
		}
		err = PublishJSON(e.Reply(), rpc)
		if err != nil {
			t.Error(err)
		}
		return nil
	}, nodeUUID)
	if err != nil {
		t.Error(err)
	}

	channelEvent := &xctrl.ChannelEvent{
		NodeUuid: nodeUUID,
	}

	channel := &Channel{
		ChannelEvent: channelEvent,
		CtrlUuid:     UUID(),
	}

	req := &xctrl.CallcenterRequest{
		Uuid: UUID(),
		Name: "test_call_center",
	}

	res := channel.Callcenter(req)

	if res.Code != 200 {
		t.Error(res)
	}
}

func TestChannel_Conference(t *testing.T) {
	// subject := "cn.xswitch.ctrl"
	//获取nats地址
	natsURL := os.Getenv("NATS_ADDRESS")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	//初始化 ctrl
	err := Init(true, natsURL)
	if err != nil {
		t.Error(err)
	}

	nodeUUID := "test.node-uuid-conference"

	//订阅主题
	_, err = Subscribe("cn.xswitch.node."+nodeUUID, func(c context.Context, e nats.Event) error {
		var request Request
		err := json.Unmarshal(e.Message().Body, &request)
		if err != nil {
			t.Error(err)
		}

		response := &xctrl.Response{
			Code:     200,
			Message:  "OK",
			NodeUuid: nodeUUID,
			Uuid:     "test-uuid",
		}

		rpc := &Response{
			Version: "2.0",
			ID:      request.ID,
			Result:  ToRawMessage(response),
		}
		err = PublishJSON(e.Reply(), rpc)
		if err != nil {
			t.Error(err)
		}
		return nil
	}, nodeUUID)
	if err != nil {
		t.Error(err)
	}

	channelEvent := &xctrl.ChannelEvent{
		NodeUuid: nodeUUID,
	}

	channel := &Channel{
		ChannelEvent: channelEvent,
		CtrlUuid:     UUID(),
	}

	req := &xctrl.ConferenceRequest{
		Uuid:    UUID(),
		Name:    "test_call_center",
		Profile: "example",
		Flags:   []string{"mute", "vmute", "deaf", "moderator", "mintwo"},
	}

	res := channel.Conference(req)

	if res.Code != 200 {
		t.Error(res)
	}

}

func TestChannel_AI(t *testing.T) {
	// subject := "cn.xswitch.ctrl"
	//获取nats地址
	natsURL := os.Getenv("NATS_ADDRESS")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	//初始化 ctrl
	err := Init(true, natsURL)
	if err != nil {
		t.Error(err)
	}

	nodeUUID := "test.node-uuid-ai"

	//订阅主题
	_, err = Subscribe("cn.xswitch.node."+nodeUUID, func(c context.Context, e nats.Event) error {
		var request Request
		err := json.Unmarshal(e.Message().Body, &request)
		if err != nil {
			t.Error(err)
		}
		response := &xctrl.Response{
			Code:     200,
			Message:  "OK",
			NodeUuid: nodeUUID,
			Uuid:     "test-uuid",
		}
		rpc := &Response{
			Version: "2.0",
			ID:      request.ID,
			Result:  ToRawMessage(response),
		}
		err = PublishJSON(e.Reply(), rpc)
		if err != nil {
			t.Error(err)
		}
		return nil
	}, nodeUUID)
	if err != nil {
		t.Error(err)
	}

	channelEvent := &xctrl.ChannelEvent{
		NodeUuid: nodeUUID,
	}

	channel := &Channel{
		ChannelEvent: channelEvent,
		CtrlUuid:     UUID(),
	}

	req := &xctrl.AIRequest{
		Uuid: UUID(),
		Url:  "http://localhost:3000",
		Data: map[string]string{
			"var1": "value1",
			"var2": "value2",
		},
	}

	res := channel.AI(req)

	if res.Code != 200 {
		t.Error(res)
	}

}

func TestChannel_HttAPI(t *testing.T) {
	// subject := "cn.xswitch.ctrl"
	//获取nats地址
	natsURL := os.Getenv("NATS_ADDRESS")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	//初始化 ctrl
	err := Init(true, natsURL)
	if err != nil {
		t.Error(err)
	}

	nodeUUID := "test.node-uuid-httapi"

	//订阅主题
	_, err = Subscribe("cn.xswitch.node."+nodeUUID, func(c context.Context, e nats.Event) error {
		var request Request
		err := json.Unmarshal(e.Message().Body, &request)
		if err != nil {
			t.Error(err)
		}
		response := &xctrl.Response{
			Code:     200,
			Message:  "OK",
			NodeUuid: nodeUUID,
			Uuid:     "test-uuid",
		}
		rpc := &Response{
			Version: "2.0",
			ID:      request.ID,
			Result:  ToRawMessage(response),
		}
		err = PublishJSON(e.Reply(), rpc)
		if err != nil {
			t.Error(err)
		}
		return nil
	}, nodeUUID)
	if err != nil {
		t.Error(err)
	}

	channelEvent := &xctrl.ChannelEvent{
		NodeUuid: nodeUUID,
	}

	channel := &Channel{
		ChannelEvent: channelEvent,
		CtrlUuid:     UUID(),
	}

	req := &xctrl.HttAPIRequest{
		Uuid: UUID(),
		Url:  "http://localhost:3000",
		Data: map[string]string{
			"var1": "value1",
			"var2": "value-httai",
		},
	}

	res := channel.HttAPI(req)

	if res.Code != 200 {
		t.Error(res)
	}

}

func TestConferenceInfo(t *testing.T) {
	// subject := "cn.xswitch.ctrl"
	//获取nats地址
	natsURL := os.Getenv("NATS_ADDRESS")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	//初始化 ctrl
	err := Init(true, natsURL)
	if err != nil {
		t.Error(err)
	}
	nodeUUID := "test.node-uuid.conferenceInfo"
	channel := &Channel{
		CtrlUuid:     UUID(),
		ChannelEvent: &xctrl.ChannelEvent{},
	}
	channel.NodeUuid = nodeUUID
	//订阅主题
	_, err = Subscribe("cn.xswitch.node."+nodeUUID, func(c context.Context, e nats.Event) error {
		var request Request
		err := json.Unmarshal(e.Message().Body, &request)
		if err != nil {
			t.Error(err)
		}
		response := &xctrl.Response{
			Code:     200,
			Message:  "OK",
			NodeUuid: nodeUUID,
			Uuid:     "test-uuid",
		}
		rpc := &Response{
			Version: "2.0",
			ID:      request.ID,
			Result:  ToRawMessage(response),
		}
		err = PublishJSON(e.Reply(), rpc)
		if err != nil {
			t.Error(err)
		}
		return nil
	}, nodeUUID)
	if err != nil {
		t.Error()
	}
	data := xctrl.ConferenceInfoRequestDataData{
		ConferenceName: "ConferenceName",
		ShowMembers:    true,
		MemberFilters: &structpb.ListValue{
			Values: []*structpb.Value{
				{
					Kind: &structpb.Value_StructValue{
						StructValue: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"role-id": &structpb.Value{
									Kind: &structpb.Value_StringValue{
										StringValue: "3",
									},
								},
								"target": &structpb.Value{
									Kind: &structpb.Value_StringValue{
										StringValue: "moderator",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	req := &xctrl.ConferenceInfoRequest{
		CtrlUuid: channel.CtrlUuid,
		Data: &xctrl.ConferenceInfoRequestData{
			Command: "conferenceInfo",
			Data:    &data,
		},
	}
	response, err := Service().ConferenceInfo(context.Background(), req, WithAddress(nodeUUID))

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response)

}

func TestLua(t *testing.T) {
	// subject := "cn.xswitch.ctrl"
	//获取nats地址
	natsURL := os.Getenv("NATS_ADDRESS")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	//初始化 ctrl
	err := Init(true, natsURL)
	if err != nil {
		t.Error(err)
	}
	nodeUUID := "test.node-uuid.lua"
	channel := &Channel{
		CtrlUuid:     UUID(),
		ChannelEvent: &xctrl.ChannelEvent{},
	}
	channel.NodeUuid = nodeUUID
	//订阅主题
	_, err = Subscribe("cn.xswitch.node."+nodeUUID, func(c context.Context, e nats.Event) error {
		var request Request
		err := json.Unmarshal(e.Message().Body, &request)
		if err != nil {
			t.Error(err)
		}
		response := &xctrl.Response{
			Code:     200,
			Message:  "OK",
			NodeUuid: nodeUUID,
			Uuid:     "test-uuid",
		}
		rpc := &Response{
			Version: "2.0",
			ID:      request.ID,
			Result:  ToRawMessage(response),
		}
		err = PublishJSON(e.Reply(), rpc)
		if err != nil {
			t.Error(err)
		}
		return nil
	}, nodeUUID)
	if err != nil {
		t.Error(err)
	}

	req := &xctrl.LuaRequest{
		Uuid:   uuid.New().String(),
		Script: "file.lua",
	}
	response, err := Service().Lua(context.Background(), req, channel.NodeAddress())

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)

}
