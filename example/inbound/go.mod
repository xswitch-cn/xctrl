module github.com/xswitch-cn/xctrl/example/inbound

go 1.16

replace github.com/xswitch-cn/xctrl => ../../

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/xswitch-cn/xctrl v1.0.0
	github.com/google/uuid v1.3.0
	github.com/sirupsen/logrus v1.7.0
)
