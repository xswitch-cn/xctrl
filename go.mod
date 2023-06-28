module git.xswitch.cn/xswitch/xctrl

go 1.16

// replace git.xswitch.cn/xswitch/proto => ../proto

require (
	git.xswitch.cn/xswitch/proto v0.1.0
	github.com/google/uuid v1.3.0
	github.com/nats-io/nats-server/v2 v2.9.18 // indirect
	github.com/nats-io/nats.go v1.27.0
	github.com/sirupsen/logrus v1.7.0
)
