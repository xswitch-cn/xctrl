# 从1.0到1.1迁移

## 1.0到1.1的主要变化

1.0到1.1的主要变化是接口变化，更规范。

### 回调机制

回调机制更新。原来是在`ctrl.Init`中传回调接口，1.1改为在`ctrl.EnableXX()`中传回调。如：

1.0：

```go
ctrl.Init(new(Handler), traceNats, "cn.xswitch.ctrl", natsAddress)
ctrl.EnableApp(Handler)
```

2.0：

```go
ctrl.Init(new(traceNats, natsAddress)
ctrl.EnableApp(new(Handler), "cn.xswitch.ctrl", "queue")
```

### Handler接口分离

以前，`ctrl.Handler`接口需要实现以下几个函数

- `App`
- `Event`
- `Request`
- `Result`

前接口分为不同的Handler，如：

- `ctrl.AppHandler`
- `ctrl.EventHandler`
- `ctrl.RequestHandler`

删除接Result，如有需要可以在`ctrl.AppHandler`中的`Event`回调中获取。

### channel接口

- 自动生成部分channel接口
- 删除了一部分channel接口，换成使用自动生成的部分
- 原来的channel接口变为加`0`后缀的接口，如果还存在可以继续使用，参数略有变化
- 推荐使用自动生成的接口，更规范，变化少。

原接口：

```go
channel.Answer()
```

新接口：

```go
channel.Answer0()
```

或

```go
channel.Answer(*xctrl.AnswerRequest)
```

### 部分参数调整

- 原来`Answer(*xctrl.Request)`改为`Answer(*xctrl.AnswerRequest)`
- 原来`NativeAPI(*xctrl.Request)`改为`NativeAPI(*xctrl.NativeAPIRequest)`
- 原来`NativeApp(*xctrl.Request)`不变
