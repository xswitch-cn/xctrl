# xctrl - XSwitch XCC API Go 语言 SDK

小樱桃在用的Go语言SDK 1.1，API更规范，推荐升级。[旧的1.0版代码和文档见这里](https://git.xswitch.cn/xswitch/xctrl/src/branch/v1.0)。

- [SDK使用和开发文档](docs/README.md)
- 协议参考文档参见：<https://git.xswitch.cn/xswitch/xctrl/src/branch/master/proto/doc>
- 示例：<https://git.xswitch.cn/xswitch/xcc-examples/src/branch/master/go>
- XCC API文档：<https://docs.xswitch.cn/xcc-api/>

目录结构：

- ctrl：节点管理
- proto：Google Protocol Buffer协议描述
- tboy 是一个冒牌的的FreeSWITCH，用于测试
- xctrl：xctrl Go语言SDK生成器，参考自Go Micro框架
- consitent 多节点一致性hash管理器，用于多个FreeSWITCH的hash节点获取，文档参见[hash节点文档](consitent.md)

更多文档参见[proto/doc](proto/doc)。

## 使用和开发

1. 克隆该项目到本地：

```shell
git clone https://git.xswitch.cn/xswitch/xctrl.git
cd xctrl
```

2. Protocol Buffers编译器（protoc）

```shell
brew install protobuf
```

3. 安装protoc-gen-doc依赖：

- 推荐方式：

```shell
make setup  
```

- 手动安装: 

```shell
go install github.com/chuanlinzhang/protoc-gen-doc/cmd/protoc-gen-doc@v0.0.2
```

4. 根据需要生成相应语言的代码:

- 生成Go代码

```shell
make proto
```

---

- 生成Java代码

```shell
make java
```
---

- 生成HTML文档

```shell
make doc-html
```
---

- 生成Markdown文档

```shell
make doc-md
```

## 测试

```shell
go run main.go
make test
```


