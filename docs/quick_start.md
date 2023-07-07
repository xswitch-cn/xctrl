# XSwitch 调用方法说明文档

Go 语言中调用 freeswitch 是通过`mod_xcc`模块通过NATS进行通信实现的，下文总结了NATS的Topic， 以及各种Go的调用方式

## Topic 列表

假设场景中有:

- 话机 : `P1`,`P2`,`P3`
- XSwitch : `X1`

当`P1` 拨打 `P2`:

- `P1` -> `X1` -> `P2`

### cn.xswitch.ctrl


### cn.xswitch.node

### cn.xswitch.ctrl

## 调用方法

`mod_xcc`的go语言SDK为`XCtrl`

通过`XCtrl`可以:

1. 调用`XNodeService`同步/异步服务接口
2. 调用`CManService`同步服务接口
3. 发布NATS消息

### 1 初始化`XCtrl`
示例

```go
ctrl.Init(true, nats_url)
if err != nil {
xlog.Panic("init ctrl err:", err)
}

```

### 2. 通过XCtrl进行调用

#### 2.1 通过XCtrl进行主动调用

- `Service()`: 同步调用 `xctrl.XNodeService` 接口
- `AsyncService()`: 异步调用 `xctrl.XNodeService` 接口
- `AService()`: 使用Context的异步调用 `xctrl.XNodeService` 接口

#### 2.2 主动调用CMan服务

- `CManService()`: 同步调用 `cman.CManService` 接口

#### 2.3 主动发布消息

- `Publish()`: 发布标准NATS消息
- `PublishJSON()`: 发布JSON格式的NATS消息

### 3. 监听事件

#### 3.1

- `EnableEvent()`: 开启事件监听
- `EnableRequest()`: 开启Request请求监听
- `EnableApp()`: 开启APP事件监听


