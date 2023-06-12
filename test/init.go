package test

import (
	"os"

	"git.xswitch.cn/xswitch/xctrl/ctrl"
)

const (
	testNodeUUID = "test.test-test"
)

func init() {
	natsURL := os.Getenv("NATS_ADDRESS")

	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	ctrl.Init(true, natsURL)
}
