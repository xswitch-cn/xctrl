package node_instance

import (
	"encoding/json"
	"git.xswitch.cn/xswitch/proto/go/proto/xctrl"
	"git.xswitch.cn/xswitch/xctrl/ctrl"
	"git.xswitch.cn/xswitch/xctrl/ctrl/nats"
	"os"
	"testing"
	"time"
)

var (
	isSuccessful = false
	natsURL      = os.Getenv("NATS_ADDRESS")
	conn         = nats.NewConn(nats.Addrs(natsURL), nats.Trace(true))
)

func TestNodeInstance(t *testing.T) {
	if natsURL == "" {
		t.Error()
		return
	}

	err := ctrl.Init(true, natsURL)
	if err != nil {
		return
	}
	ctrl.EnableNodeStatus("test.node")

	instance1, err := ctrl.NewCtrlInstance(true, natsURL)
	err = instance1.EnbaleNodeStatus("test.node1")

	err = conn.Connect()
	if err != nil {
		t.Error()
		return
	}
	node1 := &xctrl.Node{
		Uuid:                 "ins1-node1",
		Name:                 "ins1Node1",
		Ip:                   "192.168.1.1",
		Version:              "1.0.0",
		Rack:                 1,
		Address:              "192.168.1.1:8000",
		Uptime:               3600,
		Sessions:             5,
		SessionsMax:          20,
		SpsMax:               100,
		SpsLast:              5,
		SpsLast_5Min:         5,
		SessionsSinceStartup: 50,
		SessionPeak_5Min:     10,
		SessionPeakMax:       15,
	}
	publish("test.node", node1)

	node1.Uuid = "ins1-node2"
	node1.Name = "ins1Node2"
	publish("test.node", node1)

	node1.Uuid = "ins2-node1"
	node1.Name = "ins2Node1"
	publish("test.node1", node1)

	index := 0
	for {
		if index > 10 {
			break
		}
		if len(ctrl.GetNodeList()) == 2 && len(instance1.GetNodeList()) == 1 {
			isSuccessful = true
		}

		index++
		time.Sleep(time.Second)
	}

	if !isSuccessful {
		t.Error()
	}
}

func publish(topic string, param interface{}) {

	c := &ctrl.Request{
		Version: "2.0",
		Method:  "Event.NodeRegister",
		ID:      ctrl.ToRawMessage("1"),
		Params:  ctrl.ToRawMessage(param),
	}
	marshal, err := json.Marshal(c)
	if err != nil {
		return
	}

	conn.Publish(topic, marshal)
}
