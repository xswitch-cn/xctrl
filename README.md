# xctrl - XSwitch XCC API Go 语言 SDK

小樱桃在用的Go语言SDK。

- 使用文档参见：https://git.xswitch.cn/xswitch/xctrl/src/branch/master/core
- 协议参考文档参见：https://git.xswitch.cn/xswitch/xctrl/src/branch/master/core/proto/doc
- 示例：https://git.xswitch.cn/xswitch/xcc-examples/src/branch/master/go


### 使用

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
