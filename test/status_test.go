package api

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"git.xswitch.cn/xswitch/xctrl/ctrl"
	"git.xswitch.cn/xswitch/xctrl/ctrl/nats"
	"git.xswitch.cn/xswitch/xctrl/proto/xctrl"
	"git.xswitch.cn/xswitch/xctrl/xctrl/client"
)

const (
	testNodeUUID = "test.test-test"
)

type Handler struct {
}

func (h *Handler) Request(ctx context.Context, subject string, reply string, req *ctrl.Request) {
}
func (h *Handler) App(ctx context.Context, subject string, reply string, msg *ctrl.Message) {
}
func (h *Handler) Event(ctx context.Context, subject string, req *ctrl.Request) {
}
func (h *Handler) Result(ctx context.Context, subject string, result *ctrl.Result) {
}

func init() {
	natsURL := os.Getenv("NATS_ADDRESS")

	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	ctrl.Init(new(Handler), true, "cn.xswitch.ctrl."+testNodeUUID, natsURL)
}

func TestEncoding(t *testing.T) {
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
		if nativeRequest.Cmd == "status" {
			reply := &xctrl.NativeResponse{
				Code:    200,
				Message: "OK",
				Data:    "OK",
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

	response, err := ctrl.Service().NativeAPI(context.Background(), &xctrl.NativeRequest{
		Cmd: "status",
	}, client.WithAddress("cn.xswitch.node."+testNodeUUID), client.WithRequestTimeout(100*time.Millisecond))

	if err != nil {
		t.Error(err)
	}

	t.Log(response)
	// t.Error(response)
}
