module git.xswitch.cn/xswitch/xctrl/example/inbound

go 1.16

replace git.xswitch.cn/xswitch/xctrl => ../../

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	git.xswitch.cn/xswitch/xctrl v1.0.0
	github.com/google/uuid v1.3.0
	github.com/sirupsen/logrus v1.7.0
)
