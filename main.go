package main

import (
	"context"
	"time"

	"log"

	"git.xswitch.cn/xswitch/xctrl/ctrl"
	"git.xswitch.cn/xswitch/xctrl/proto/cman"
	"git.xswitch.cn/xswitch/xctrl/proto/xctrl"
)

type Logger struct {
	ctrl.Logger
}

func (t *Logger) Log(level int, v ...interface{}) {
	// if ctrl.LogLevel(level) == ctrl.LLTrace {
	// }
	log.Println(v...)
}

func (t *Logger) Logf(level int, format string, v ...interface{}) {
	// if ctrl.LogLevel(level) == ctrl.LLTrace {
	// }
	log.Printf(format, v...)
}

// simple example
func main() {
	logLevel := ctrl.LLDebug
	// logLevel = ctrl.LLTrace // uncomment this line to enable trace log
	isTrace := logLevel == ctrl.LLTrace // should enable trace log in ctrl?
	ctrl.SetLogLevel(logLevel)          // set ctrl log level
	ctrl.SetLogger(new(Logger))         // tell ctrl to use our logger
	log.Print("Hello, world!")          // the world starts from here
	// init ctrl, connect to NATS and subscribe a subject
	err := ctrl.Init(isTrace, "nats://localhost:4222")
	if err != nil {
		panic(err)
	}
	ctrl.EnableApp(new(ctrl.EmptyAppHandler), "cn.xswitch.ctrl", "ctrl")
	ctrl.EnableNodeStatus("")
	// init cman service before we can talk to cman
	ctrl.InitCManService("cn.xswitch.cman.control")

	response, err := ctrl.Service().NativeAPI(context.Background(), &xctrl.NativeAPIRequest{
		Cmd: "status",
	}, ctrl.WithAddress("cn.xswitch.node"), ctrl.WithRequestTimeout(1*time.Second))

	if err != nil {
		panic(err)
	}

	log.Printf("response: %v", response.Data)

	cListReq := &xctrl.ConferenceListRequest{
		CtrlUuid: ctrl.UUID(),
		Data: &xctrl.ConferenceListRequestData{
			Command: "conferenceInfo",
			Data: &xctrl.ConferenceListRequestDataData{
				Domain: "",
			},
		},
	}
	rsp, err := ctrl.Service().ConferenceList(context.Background(), cListReq,
		ctrl.WithAddress("cn.xswitch.node"), ctrl.WithRequestTimeout(1*time.Second))
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("%v", rsp.Data)
		for _, c := range rsp.Data {
			log.Printf("conference %s %s", c.ConferenceName, c.Domain)
		}
	}

	res, err := ctrl.CManService().GetConferenceList(context.Background(), &cman.GetConferenceListRequest{},
		ctrl.WithRequestTimeout(1*time.Second))
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("conferences %v", res.Conferences)
	}
}
