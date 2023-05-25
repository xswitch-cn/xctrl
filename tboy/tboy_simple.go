package tboy

import (
	"context"
	"encoding/json"

	"git.xswitch.cn/xswitch/xctrl/ctrl"
	"git.xswitch.cn/xswitch/xctrl/ctrl/nats"
)

type TBoySimple struct {
	*TBoy
}

// NewTBoy .
func NewSimple(node_uuid string, domain string, fn ...OptionFn) *TBoySimple {
	boy := &TBoySimple{
		TBoy: &TBoy{
			Options: &Options{},
		},
	}
	boy.Init()
	boy.SetUUID(node_uuid)
	boy.SetDomain(domain)
	for _, o := range fn {
		o(boy.Options)
	}

	log.Infof("new boy created uuid=%s domain=%s", boy.NodeUUID, boy.Domain)
	return boy
}

func (boy *TBoySimple) AddChannel(key string, channel *FakeChannel) {
	Channels[key] = channel
}

func (boy *TBoySimple) EventCallback(ctx context.Context, ev nats.Event) error {
	log.Info(ev.Topic(), string(ev.Message().Body))

	var msg ctrl.Message
	err := json.Unmarshal(ev.Message().Body, &msg)

	if err != nil {
		log.Error("parse error", ev)
		return err
	}

	boy.App(ctx, ev.Topic(), ev.Reply(), &msg)

	return nil
}

// App .
func (boy *TBoySimple) App(ctx context.Context, topic string, reply string, msg *ctrl.Message) {
	log.Infof("Handle %s %s", topic, msg.Method)

	switch msg.Method {
	default:
		// handle messages in TBoy
		boy.TBoy.App(ctx, topic, reply, msg)
	}

}

// Event .
func (boy *TBoySimple) Event(ctx context.Context, topic string, message *ctrl.Request) {
}

// Request .
func (boy *TBoySimple) Request(ctx context.Context, topic string, reply string, request *ctrl.Request) {
}

// Result 异步请求结果
func (boy *TBoySimple) Result(ctx context.Context, topic string, message *ctrl.Result) {}
