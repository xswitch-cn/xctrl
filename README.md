# xctrl - XSwitch XCC API Go 语言 SDK

小樱桃在用的Go语言SDK 1.1，API更规范，推荐升级。[旧的1.0版代码和文档见这里](https://git.xswitch.cn/xswitch/xctrl/src/branch/v1.0)。

- 依赖 <https://git.xswitch.cn/xswitch/proto> 。
- [SDK使用和开发文档](docs/README.md)
- 协议参考文档参见：<https://git.xswitch.cn/xswitch/proto/src/branch/main/docs>
- 示例：<https://git.xswitch.cn/xswitch/xcc-examples/src/branch/master/go>
- XCC API文档：<https://docs.xswitch.cn/xcc-api/>

目录结构：

- ctrl：节点管理
- tboy 是一个冒牌的的FreeSWITCH，用于测试
- consitent 多节点一致性hash管理器，用于多个FreeSWITCH的hash节点获取，文档参见[hash节点文档](docs/consitent.md)

## 使用

```sh
go mod init main
go get git.xswitch.cn/xswitch/proto/go/proto/xctrl
go get git.xswitch.cn/xswitch/xctrl
# 创建 main.go，并运行
go run main.go
```

## 开发

1. 克隆该项目到本地：

```shell
git clone https://git.xswitch.cn/xswitch/xctrl.git
cd xctrl
```

## 测试

```shell
go run main.go
make test
```
