package main

import (
	"context"
	"time"

	"log"

	"github.com/xswitch-cn/proto/go/proto/cman"
	"github.com/xswitch-cn/proto/go/proto/xctrl"
	"github.com/xswitch-cn/xctrl/ctrl"
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
	logLevel = ctrl.LLTrace             // uncomment this line to enable trace log
	isTrace := logLevel == ctrl.LLTrace // should enable trace log in ctrl?
	ctrl.SetLogLevel(logLevel)          // set ctrl log level
	ctrl.SetLogger(new(Logger))         // tell ctrl to use our logger
	log.Print("Hello, world!")          // the world starts from here
	// init ctrl, connect to NATS and subscribe a subject
	err := ctrl.Init(isTrace, "nats://localhost:4222")
	// err := ctrl.Init(isTrace, "nats://user:pass@localhost:4222")
	if err != nil {
		panic(err)
	}
	tenant := ""
	// tenant = "cherry"
	prefix := ""
	if tenant != "" {
		ctrl.SetFromPrefix("from-")
		ctrl.SetToPrefix("to-")
		prefix = "to-" + tenant + "."
	}

	// ctrl.EnableApp(new(ctrl.EmptyAppHandler), "cn.xswitch.ctrl", "ctrl")
	ctrl.EnableNodeStatus("")
	// init cman service before we can talk to cman
	ctrl.InitCManService(prefix + "cn.xswitch.cman.control")

	response, err := ctrl.Service().NativeAPI(context.Background(), &xctrl.NativeAPIRequest{
		Cmd: "status",
	}, ctrl.WithTenantAddress("cherry", "cn.xswitch.node"), ctrl.WithRequestTimeout(1*time.Second))

	if err != nil {
		panic(err)
	}

	log.Printf("response: %v", response.Data)
	_, err = ctrl.Service().NativeAPI(context.Background(), &xctrl.NativeAPIRequest{
		Cmd:  "log",
		Args: "INFO xctrl test log",
	}, ctrl.WithTenantAddress("cherry", "cn.xswitch.node"), ctrl.WithAsync())

	if err != nil {
		panic(err)
	}

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
		ctrl.WithTenantAddress("cherry", "cn.xswitch.node"), ctrl.WithRequestTimeout(1*time.Second))
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("%v", rsp.Data)
		for _, c := range rsp.Data {
			log.Printf("conference %s %s", c.ConferenceName, c.Domain)
		}
	}

	option := ctrl.WithTenantAddress("cherry", "cn.xswitch.cman.control")
	res, err := ctrl.CManService().GetConferenceList(context.Background(), &cman.GetConferenceListRequest{},
		ctrl.WithRequestTimeout(1*time.Second), option)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("conferences %v", res.Conferences)
	}
}
