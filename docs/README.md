# xctrl 文档

本文档是xctrl 1.4文档。

需要注意，xctrl并不使用Protocol Buffer以及gRPC，而是使用了Protocol Buffer的协议描述，方便生成跨语言的客户端户。xctrl与XSwitch之间使用JSON数据格式，并使用JSON-RPC封装。

本SDK依赖`google.golang.org/protobuf/`包，本SDK有两种使用方式：

- 直接使用`ctrl`和`xctrl`包，集成了协议结构体（`proto/xctrl/xctrl.pb.go`）及NATS消息收发（`proto/xctrl/xctrl.pb.xctrl.go`）。
- 仅使用生成的Go结构体，如仅使用`xctrl.pb.go`，而不使用`xctrl.pb.xctrl.go`，可以直接将这两个文件复制到你的项目中。

下面是相应的文档。

- [1.0向1.1迁移](migration.md)
- [SDK使用](sdk.md)
- [ctrl](ctrl.md)
    - [多租户](ctrl_tenant.md)
    - [多实例](ctrl_instance.md)
- [channel](channel.md)
- [SDK](sdk.md)
- [DTMF](dtmf.md)
- [bus](bus.md)
- [开发](dev.md)
- [FSDS](fsds.md)
