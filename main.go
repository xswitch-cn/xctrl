package main

import (
	"context"
	"time"

	"git.xswitch.cn/xswitch/xctrl/core/ctrl"
	"git.xswitch.cn/xswitch/xctrl/core/proto/cman"
	"git.xswitch.cn/xswitch/xctrl/xctrl/client"
	"git.xswitch.cn/xswitch/xctrl/xctrl/util/log"
)

type Handler struct {
}

func (h *Handler) Request(context.Context, string, string, *ctrl.Request) {
}
func (h *Handler) App(context.Context, string, string, *ctrl.Message) {
}
func (h *Handler) Event(context.Context, string, *ctrl.Request) {
}
func (h *Handler) Result(context.Context, string, *ctrl.Result) {
}

// simple example
func main() {
	log.SetLevel(log.LevelDebug)
	// log.SetLogger(...)
	log.Info("Hello, world!")

	ctrl.Init(new(Handler), true, "cn.xswitch.ctrl", "nats://localhost:4222")
	ctrl.InitCManService("cn.xswitch.ctrl.cman")

	ctrl.CManService().GetConferenceList(context.Background(), &cman.GetConferenceListRequest{}, ctrl.WithAddress("cn.xswitch.ctrl.cman"), client.WithRequestTimeout(1*time.Second))

}
