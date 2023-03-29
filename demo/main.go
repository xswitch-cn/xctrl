package main

import (
	"fmt"
	xlog "git.xswitch.cn/xswitch/xctrl/log"
)

func main() {
	logCallback := func(Level xlog.LevelLog, logText interface{}) error {
		level, _ := Level.MarshalText()
		switch Level {
		case xlog.TraceLevel:
			fmt.Println("Level:", level, "-log:", logText)
		case xlog.DebugLevel:
			fmt.Println("Level:", level, "-log:", logText)
		case xlog.WarnLevel:
			fmt.Println("Level:", level, "-log:", logText)
		case xlog.InfoLevel:
			fmt.Println("Level:", level, "-log:", logText)
		case xlog.ErrorLevel:
			fmt.Println("Level:", level, "-log:", logText)
		case xlog.FatalLevel:
			fmt.Println("Level:", level, "-log:", logText)
		default:
			fmt.Println("Level:", level, "-log:", logText)
		}
		return nil
	}
	xlog.Init(xlog.InfoLevel, logCallback)
	xlog.TestLog("demo test log")
}
