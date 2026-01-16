package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/xswitch-cn/proto/go/proto/xctrl"
	"github.com/xswitch-cn/xctrl/ctrl"

	"github.com/xswitch-cn/xctrl/tboy"
)

var (
	domain      = "test.test"
	ctrlSubject = "cn.xswitch.ctrl"
)

func PubStart(node_uuid string) {
	channelEvent := xctrl.ChannelEvent{
		NodeUuid:    node_uuid,
		Uuid:        uuid.New().String(),
		Domain:      domain,
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
	ctrl.Publish(ctrlSubject, req_str)
}

func main() {
	var log = logrus.New()
	log.SetReportCaller(true)
	natsAddress := os.Getenv("NATS_ADDRESS")
	if natsAddress == "" {
		natsAddress = "nats://127.0.0.1:4222"
	}
	log.Infof("connecting to nats: %s", natsAddress)
	err := ctrl.Init(true, natsAddress)
	if err != nil {
		log.Fatal("ctrl init failed: ", err)
	}
	node_uuid := ctrl.UUID() // use ctrl uuid as node uuid
	boy := tboy.NewSimple(node_uuid, domain, tboy.OptionPeerAnswer(true), tboy.OptionActualPlay(false))

	// subscribe to cn.xswitch.node and cn.xswitch.node.$node_uuid
	ctrl.EnableApp(boy, "cn.xswitch.node", "q")

	PubStart(node_uuid)
	time.Sleep(time.Second * 30)
}
