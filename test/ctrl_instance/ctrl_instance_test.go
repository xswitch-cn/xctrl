package ctrl_instance

import (
	"github.com/xswitch-cn/xctrl/ctrl"
	"github.com/xswitch-cn/xctrl/ctrl/nats"
	"os"
	"testing"
	"time"
)

var (
	isSuccessful = true
	natsURL      = os.Getenv("NATS_ADDRESS")
	msg          = "{\"Hello\":\"World\"}"
	msg1         = "{\"Hello\":\"World1\"}"
	msg2         = "{\"Hello\":\"World2\"}"
)

type CtrlInstanceEvent struct{}

func (e CtrlInstanceEvent) Event(req *ctrl.Request, natsEvent nats.Event) {
	if string(natsEvent.Message().Body) != msg {
		isSuccessful = false
	}
}

type CtrlInstanceEvent1 struct{}

func (e CtrlInstanceEvent1) Event(req *ctrl.Request, natsEvent nats.Event) {
	if string(natsEvent.Message().Body) != msg1 {
		isSuccessful = false
	}
}

type CtrlInstanceEvent2 struct{}

func (e CtrlInstanceEvent2) Event(req *ctrl.Request, natsEvent nats.Event) {
	if string(natsEvent.Message().Body) != msg2 {
		isSuccessful = false
	}
}

func TestCtrlInstance(t *testing.T) {
	if natsURL == "" {
		t.Error()
		return
	}

	err := ctrl.Init(true, natsURL)
	if err != nil {
		return
	}
	ctrl.EnableEvent(new(CtrlInstanceEvent), "test.test", "")
	ctrl.EnableNodeStatus("test.node")
	list := ctrl.GetNodeList()
	t.Log(list)

	instance1, err := ctrl.NewCtrlInstance(true, natsURL)
	if err != nil {
		return
	}
	instance1.EnableEvent(new(CtrlInstanceEvent1), "test.test1", "")

	err = instance1.EnbaleNodeStatus("test.node1")
	if err != nil {
		return
	}
	t.Log(instance1.GetNodeList())

	instance2, err := ctrl.NewCtrlInstance(true, natsURL)
	if err != nil {
		return
	}
	instance2.EnableEvent(new(CtrlInstanceEvent2), "test.test2", "")

	conn := nats.NewConn(nats.Addrs(natsURL), nats.Trace(true))
	err = conn.Connect()
	if err != nil {
		t.Error()
		return
	}
	conn.Publish("test.test", []byte(msg))
	conn.Publish("test.test1", []byte(msg1))
	conn.Publish("test.test2", []byte(msg2))

	time.Sleep(3 * time.Second)
	if !isSuccessful {
		t.Error()
	}
}
