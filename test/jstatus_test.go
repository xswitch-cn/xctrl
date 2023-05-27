package test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"git.xswitch.cn/xswitch/xctrl/ctrl"
	"git.xswitch.cn/xswitch/xctrl/ctrl/nats"
	"git.xswitch.cn/xswitch/xctrl/proto/xctrl"
	"git.xswitch.cn/xswitch/xctrl/xctrl/client"
	"git.xswitch.cn/xswitch/xctrl/xctrl/util/log"
)

func service(c context.Context, e nats.Event) error {
	msg := e.Message()
	req := &ctrl.Request{}
	err := json.Unmarshal(msg.Body, &req)
	if err != nil {
		return err
	}
	request := &xctrl.JStatusRequest{}
	err = json.Unmarshal(*req.Params, &request)
	if err != nil {
		return err
	}
	if request.Data == nil {
		log.Error("request data is nil")
		return err
	}
	log.Info(request.Data.Command)
	if request.Data.Command == "status" {
		reply := &xctrl.JStatusResponse{
			Code:    200,
			Message: "OK",
			Data: &xctrl.JStatusResponseData{
				SystemStatus: "running",
				Version:      "1.0",
			},
		}
		result := &ctrl.Result{
			Version: "2.0",
			ID:      req.ID,
			Result:  ctrl.ToRawMessage(reply),
		}
		ctrl.Publish(e.Reply(), *ctrl.ToRawMessage(result))
	}
	return nil
}

func TestJStatus(t *testing.T) {
	ctrl.Subscribe("cn.xswitch.node."+testNodeUUID, service, "")

	response, err := ctrl.Service().JStatus(context.Background(), &xctrl.JStatusRequest{
		CtrlUuid: "cn.xswitch.ctrl." + testNodeUUID,
		Data: &xctrl.JStatusRequest_JStatusData{
			Command: "status",
		},
	}, ctrl.WithAddress("cn.xswitch.node."+testNodeUUID), client.WithRequestTimeout(100*time.Millisecond))

	if err != nil {
		t.Error(err)
	}

	t.Log("response", response)

	if response.Data.SystemStatus != "running" {
		t.Error("status error")
	}
}
