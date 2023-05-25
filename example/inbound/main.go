package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"git.xswitch.cn/xswitch/xctrl/ctrl"
	"git.xswitch.cn/xswitch/xctrl/proto/xctrl"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"git.xswitch.cn/xswitch/xctrl/tboy"
)

var (
	node_topic = "test.simple-test"
	domain     = "test.test"
	node_uuid  = "test.simple-test.simple"
	subject    = "cn.xswitch.ctrl"
)

func PubStart() {
	channelEvent := xctrl.ChannelEvent{
		NodeUuid:    node_uuid,
		Uuid:        uuid.New().String(),
		Direction:   "outbound",
		State:       "START",
		CidName:     "TEST",
		CidNumber:   "1001",
		DestNumber:  "1002",
		AnswerEpoch: uint32(time.Now().Unix()),
		Answered:    true,
		CreateEpoch: uint32(time.Now().Unix()),
		Params:      map[string]string{},
	}
	channelEvent.Params["xcc_session"] = uuid.New().String()
	channel := &tboy.FakeChannel{
		CtrlUuid: "",
		Data:     &channelEvent,
	}
	tboy.CacheChannel(channelEvent.Uuid, channel)
	event_req := ctrl.Request{
		Version: "2.0",
		Method:  "Event.Channel",
		Params:  ctrl.ToRawMessage(channelEvent),
	}

	req_str, _ := json.MarshalIndent(event_req, "", "  ")
	ctrl.Publish(subject, req_str)
}

func main() {
	var log = logrus.New()
	log.SetReportCaller(true)
	boy := tboy.NewSimple(node_uuid, domain, tboy.OptionPeerAnswer(true), tboy.OptionActualPlay(false))

	natsAddress := os.Getenv("NATS_ADDRESS")

	if natsAddress == "" {
		natsAddress = "nats://127.0.0.1:4222"
	}
	log.Infof("connecting to nats: %s", natsAddress)
	err := ctrl.Init(boy, true, "cn.xswitch.ctrl", natsAddress)
	fmt.Println(err)
	if err != nil {
		log.Fatal("ctrl init failed: ", err)
	}

	ctrl.Subscribe("cn.xswitch.node", boy.EventCallback, "node")
	ctrl.Subscribe("cn.xswitch.node.test", boy.EventCallback, "node")
	ctrl.Subscribe("cn.xswitch.node."+node_topic, boy.EventCallback, "node")
	ctrl.Subscribe("cn.xswitch.node."+node_uuid, boy.EventCallback, "")

	PubStart()
	time.Sleep(time.Second * 30)

}
