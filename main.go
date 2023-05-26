package main

import (
	"context"
	"time"

	"git.xswitch.cn/xswitch/xctrl/ctrl"
	"git.xswitch.cn/xswitch/xctrl/proto/cman"
	"git.xswitch.cn/xswitch/xctrl/proto/xctrl"
	"git.xswitch.cn/xswitch/xctrl/xctrl/client"
	"git.xswitch.cn/xswitch/xctrl/xctrl/util/log"
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

// simple example
func main() {
	log.SetLevel(log.LevelDebug)
	// log.SetLogger(...)
	log.Info("Hello, world!")

	ctrl.Init(new(Handler), true, "cn.xswitch.ctrl", "nats://localhost:4222")
	ctrl.InitCManService("cn.xswitch.ctrl.cman")

	response, err := ctrl.Service().NativeAPI(context.Background(), &xctrl.NativeRequest{
		Cmd: "status",
	}, ctrl.WithAddress("cn.xswitch.node"), client.WithRequestTimeout(1*time.Second))

	if err != nil {
		panic(err)
	}

	log.Infof("response: %v", response.Data)

	res, err := ctrl.CManService().GetConferenceList(context.Background(), &cman.GetConferenceListRequest{},
		ctrl.WithAddress("cn.xswitch.cman.control"), client.WithRequestTimeout(1*time.Second))
	if err != nil {
		log.Error(err)
	} else {
		log.Info("conferences", res.Conferences)
	}
}
