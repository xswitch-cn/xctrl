package test

import (
	"os"

	"github.com/xswitch-cn/xctrl/ctrl"
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
