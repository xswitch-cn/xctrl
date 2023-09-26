package node_instance

import (
	"context"
	"git.xswitch.cn/xswitch/xctrl/ctrl"
	"git.xswitch.cn/xswitch/xctrl/ctrl/nats"
	"testing"
	"time"
)

var successful = false

func TestTenantAddress(t *testing.T) {
	tenant := "foo"
	nodeUUID := "test.tenant"

	instance, err := ctrl.NewCtrlInstance(true, "nats://192.168.3.235:4222")
	if err != nil {
		return
	}

	instance.Subscribe(instance.TenantNodeAddress(tenant, nodeUUID), func(ctx context.Context, event nats.Event) error {

		s := string(event.Message().Body)
		if "hello" == s {
			successful = true
		}
		return nil
	}, "")

	instance.Publish(instance.TenantNodeAddress(tenant, nodeUUID), []byte("hello"))
	if err != nil {
		t.Error("Publish error")
	}
	time.Sleep(2 * time.Second)

	if !successful {
		t.Error("Test failed")
	}
}
