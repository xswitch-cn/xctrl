package test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"git.xswitch.cn/xswitch/proto/go/proto/xctrl"
	"git.xswitch.cn/xswitch/proto/xctrl/client"
	"git.xswitch.cn/xswitch/xctrl/ctrl"
	"git.xswitch.cn/xswitch/xctrl/ctrl/nats"
)

func TestNativeAPI(t *testing.T) {
	ctrl.Subscribe("cn.xswitch.node."+testNodeUUID, func(c context.Context, e nats.Event) error {
		msg := e.Message()
		req := &ctrl.Request{}
		err := json.Unmarshal(msg.Body, &req)
		if err != nil {
			t.Error(err)
		}
		nativeRequest := &xctrl.NativeRequest{}
		err = json.Unmarshal(*req.Params, &nativeRequest)
		if err != nil {
			t.Error(err)
		}
		t.Log(nativeRequest.Cmd)
		if nativeRequest.Cmd == "ping" {
			reply := &xctrl.NativeResponse{
				Code:    200,
				Message: "OK",
				Data:    "pong",
			}
			result := &ctrl.Result{
				Version: "2.0",
				ID:      req.ID,
				Result:  ctrl.ToRawMessage(reply),
			}
			ctrl.Publish(e.Reply(), *ctrl.ToRawMessage(result))
		}
		return nil
	}, "")

	response, err := ctrl.Service().NativeAPI(context.Background(), &xctrl.NativeAPIRequest{
		Cmd: "ping",
	}, client.WithAddress("cn.xswitch.node."+testNodeUUID), client.WithRequestTimeout(100*time.Millisecond))

	if err != nil {
		t.Error(err)
	}

	t.Log(response)

	if response.Data != "pong" {
		t.Error("response data is not pong")
	}
}
