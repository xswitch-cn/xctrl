package ctrl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"
	"testing"
	"time"

	"git.xswitch.cn/xswitch/xctrl/core/ctrl/nats"
	"git.xswitch.cn/xswitch/xctrl/core/proto/xctrl"
)

func TestPlayWithTimeout(t *testing.T) {
	subject := "cn.xswitch.ctrl"
	url := os.Getenv("NATS_ADDRESS")

	if url == "" {
		url = "nats://localhost:4222"
	}

	err := Init(nil, true, subject, url)
	if err != nil {
		t.Error(err)
	}

	node_uuid := "test.node-uuid"

	channel := &Channel{
		CtrlUuid: UUID(),
	}

	channel.NodeUuid = node_uuid

	req := &xctrl.PlayRequest{
		CtrlUuid: UUID(),
		Uuid:     "test-uuid",
		Media: &xctrl.Media{
			Data: "/tmp/test.wav",
		},
	}

	res := channel.PlayWithTimeout(req, 100*time.Millisecond)

	if res.Code != 408 {
		t.Error(res)
	}

	Subscribe("cn.xswitch.node."+node_uuid, func(c context.Context, e nats.Event) error {
		var request Request
		json.Unmarshal(e.Message().Body, request)

		response := &xctrl.Response{
			Code:     200,
			Message:  "OK",
			NodeUuid: node_uuid,
			Uuid:     "test-uuid",
		}

		rpc := &Response{
			Version: "2.0",
			ID:      request.ID,
			Result:  ToRawMessage(response),
		}

		PublishJSON(e.Reply(), rpc)
		return nil
	}, node_uuid)

	res = channel.PlayWithTimeout(req, 100*time.Millisecond)

	if res.Code != 200 {
		t.Error(res)
	}

	res = channel.Play(req)

	if res.Code != 200 {
		t.Error(res)
	}
}

func TestFIFO(t *testing.T) {
	subject := "cn.xswitch.ctrl"
	//获取nats地址
	url := os.Getenv("NATS_ADDRESS")
	if url == "" {
		url = "nats://localhost:4222"
	}
	//初始化 ctrl
	err := Init(nil, true, subject, url)
	if err != nil {
		t.Error(err)
	}

	node_uuid := "test.node-uuid"

	//订阅主题
	Subscribe("cn.xswitch.node."+node_uuid, func(c context.Context, e nats.Event) error {
		var request Request
		json.Unmarshal(e.Message().Body, request)

		response := &xctrl.Response{
			Code:     200,
			Message:  "OK",
			NodeUuid: node_uuid,
			Uuid:     "test-uuid",
		}

		rpc := &Response{
			Version: "2.0",
			ID:      request.ID,
			Result:  ToRawMessage(response),
		}
		PublishJSON(e.Reply(), rpc)
		return nil
	}, node_uuid)

	channel := &Channel{
		CtrlUuid: UUID(),
	}
	channel.NodeUuid = node_uuid

	req := &xctrl.FIFORequest{
		Uuid:         UUID(),
		Name:         "test_name",
		Inout:        "out",
		WaitMusic:    "/tmp/test.wav",
		ExitAnnounce: "/tmp/test.wav",
		Priority:     0,
	}

	res := channel.FIFO(req)

	if res.Code != 408 {
		t.Error(res)
	}

}

func TestChannel_Callcenter(t *testing.T) {
	subject := "cn.xswitch.ctrl"
	//获取nats地址
	url := os.Getenv("NATS_ADDRESS")
	if url == "" {
		url = "nats://localhost:4222"
	}
	//初始化 ctrl
	err := Init(nil, true, subject, url)
	if err != nil {
		t.Error(err)
	}

	node_uuid := "test.node-uuid"

	//订阅主题
	Subscribe("cn.xswitch.node."+node_uuid, func(c context.Context, e nats.Event) error {
		var request Request
		json.Unmarshal(e.Message().Body, request)

		response := &xctrl.Response{
			Code:     200,
			Message:  "OK",
			NodeUuid: node_uuid,
			Uuid:     "test-uuid",
		}

		rpc := &Response{
			Version: "2.0",
			ID:      request.ID,
			Result:  ToRawMessage(response),
		}
		PublishJSON(e.Reply(), rpc)
		return nil
	}, node_uuid)

	channel := &Channel{
		CtrlUuid: UUID(),
	}
	channel.NodeUuid = node_uuid

	req := &xctrl.CallcenterRequest{
		Uuid: UUID(),
		Name: "test_call_center",
	}

	res := channel.Callcenter(req)

	if res.Code != 408 {
		t.Error(res)
	}

}

func TestChannel_Conference(t *testing.T) {
	subject := "cn.xswitch.ctrl"
	//获取nats地址
	url := os.Getenv("NATS_ADDRESS")
	if url == "" {
		url = "nats://localhost:4222"
	}
	//初始化 ctrl
	err := Init(nil, true, subject, url)
	if err != nil {
		t.Error(err)
	}

	node_uuid := "test.node-uuid-conference"

	//订阅主题
	Subscribe("cn.xswitch.node."+node_uuid, func(c context.Context, e nats.Event) error {
		var request Request
		json.Unmarshal(e.Message().Body, request)

		response := &xctrl.Response{
			Code:     200,
			Message:  "OK",
			NodeUuid: node_uuid,
			Uuid:     "test-uuid",
		}

		rpc := &Response{
			Version: "2.0",
			ID:      request.ID,
			Result:  ToRawMessage(response),
		}
		PublishJSON(e.Reply(), rpc)
		return nil
	}, node_uuid)

	channel := &Channel{
		CtrlUuid: UUID(),
	}
	channel.NodeUuid = node_uuid

	req := &xctrl.ConferenceRequest{
		Uuid:    UUID(),
		Name:    "test_call_center",
		Profile: "example",
		Flags:   []string{"mute", "vmute", "deaf", "moderator", "mintwo"},
	}

	res := channel.Conference(req)

	if res.Code != 408 {
		t.Error(res)
	}

}

func TestChannel_AI(t *testing.T) {
	subject := "cn.xswitch.ctrl"
	//获取nats地址
	url := os.Getenv("NATS_ADDRESS")
	if url == "" {
		url = "nats://localhost:4222"
	}
	//初始化 ctrl
	err := Init(nil, true, subject, url)
	if err != nil {
		t.Error(err)
	}

	node_uuid := "test.node-uuid-ai"

	//订阅主题
	Subscribe("cn.xswitch.node."+node_uuid, func(c context.Context, e nats.Event) error {
		var request Request
		json.Unmarshal(e.Message().Body, request)
		response := &xctrl.Response{
			Code:     200,
			Message:  "OK",
			NodeUuid: node_uuid,
			Uuid:     "test-uuid",
		}
		rpc := &Response{
			Version: "2.0",
			ID:      request.ID,
			Result:  ToRawMessage(response),
		}
		PublishJSON(e.Reply(), rpc)
		return nil
	}, node_uuid)

	channel := &Channel{
		CtrlUuid: UUID(),
	}
	channel.NodeUuid = node_uuid
	req := &xctrl.AIRequest{
		Uuid: UUID(),
		Url:  "http://localhost:3000",
		Data: map[string]string{
			"var1": "value1",
			"var2": "value2",
		},
	}

	res := channel.AI(req)

	if res.Code != 408 {
		t.Error(res)
	}

}

func TestChannel_HttAPI(t *testing.T) {
	subject := "cn.xswitch.ctrl"
	//获取nats地址
	url := os.Getenv("NATS_ADDRESS")
	if url == "" {
		url = "nats://localhost:4222"
	}
	//初始化 ctrl
	err := Init(nil, true, subject, url)
	if err != nil {
		t.Error(err)
	}

	node_uuid := "test.node-uuid-httapi"

	//订阅主题
	Subscribe("cn.xswitch.node."+node_uuid, func(c context.Context, e nats.Event) error {
		var request Request
		json.Unmarshal(e.Message().Body, request)
		response := &xctrl.Response{
			Code:     200,
			Message:  "OK",
			NodeUuid: node_uuid,
			Uuid:     "test-uuid",
		}
		rpc := &Response{
			Version: "2.0",
			ID:      request.ID,
			Result:  ToRawMessage(response),
		}
		PublishJSON(e.Reply(), rpc)
		return nil
	}, node_uuid)

	channel := &Channel{
		CtrlUuid: UUID(),
	}
	channel.NodeUuid = node_uuid
	req := &xctrl.HttAPIRequest{
		Uuid: UUID(),
		Url:  "http://localhost:3000",
		Data: map[string]string{
			"var1": "value1",
			"var2": "value-httai",
		},
	}

	res := channel.HttAPI(req)

	if res.Code != 408 {
		t.Error(res)
	}

}

func TestConferenceInfo(t *testing.T) {
	subject := "cn.xswitch.ctrl"
	//获取nats地址
	url := os.Getenv("NATS_ADDRESS")
	if url == "" {
		url = "nats://localhost:4222"
	}
	//初始化 ctrl
	err := Init(nil, true, subject, url)
	if err != nil {
		t.Error(err)
	}
	node_uuid := "test.node-uuid.conferenceInfo"
	channel := &Channel{
		CtrlUuid: UUID(),
	}
	channel.NodeUuid = node_uuid
	//订阅主题
	Subscribe("cn.xswitch.node."+node_uuid, func(c context.Context, e nats.Event) error {
		var request Request
		json.Unmarshal(e.Message().Body, request)
		response := &xctrl.Response{
			Code:     200,
			Message:  "OK",
			NodeUuid: node_uuid,
			Uuid:     "test-uuid",
		}
		rpc := &Response{
			Version: "2.0",
			ID:      request.ID,
			Result:  ToRawMessage(response),
		}
		PublishJSON(e.Reply(), rpc)
		return nil
	}, node_uuid)
	data := xctrl.ConferenceInfoRequestDataData{
		ConferenceName: "ConferenceName",
		ShowMembers:    true,
		MemberFilters: map[string]string{
			"role-id": "3",
			"target":  "moderator",
		},
	}
	req := &xctrl.ConferenceInfoRequest{
		CtrlUuid: channel.CtrlUuid,
		Data: &xctrl.ConferenceInfoRequestData{
			Command: "conferenceInfo",
			Data:    &data,
		},
	}
	response, err := Service().ConferenceInfo(context.Background(), req, channel.NodeAddress())

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response)

}

func TestLua(t *testing.T) {
	subject := "cn.xswitch.ctrl"
	//获取nats地址
	url := os.Getenv("NATS_ADDRESS")
	if url == "" {
		url = "nats://localhost:4222"
	}
	//初始化 ctrl
	err := Init(nil, true, subject, url)
	if err != nil {
		t.Error(err)
	}
	node_uuid := "test.node-uuid.lua"
	channel := &Channel{
		CtrlUuid: UUID(),
	}
	channel.NodeUuid = node_uuid
	//订阅主题
	Subscribe("cn.xswitch.node."+node_uuid, func(c context.Context, e nats.Event) error {
		var request Request
		json.Unmarshal(e.Message().Body, request)
		response := &xctrl.Response{
			Code:     200,
			Message:  "OK",
			NodeUuid: node_uuid,
			Uuid:     "test-uuid",
		}
		rpc := &Response{
			Version: "2.0",
			ID:      request.ID,
			Result:  ToRawMessage(response),
		}
		PublishJSON(e.Reply(), rpc)
		return nil
	}, node_uuid)

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
