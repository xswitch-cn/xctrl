# 迁移文档

## 1.0到1.1

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


## 1.3 到1.4

1.3 到 1.4 的主要变化是在`ctrl.AppHandler` 中 `ChannelEvent` 的 `context.Context` 可以在同一通道的多次回调之间“传递/累积”，便于在 `START` 阶段把数据放进 `ctx`，并在后续 `RINGING/ANSWERED/DESTROY` 等事件继续读取。

### ChannelEvent 返回值变更

`ctrl.AppHandler` 的 `ChannelEvent` 回调签名从“无返回值”调整为返回 `context.Context`：

v1.3：

```go
ChannelEvent(ctx context.Context, channel *ctrl.Channel)
```

v1.4：

```go
ChannelEvent(ctx context.Context, channel *ctrl.Channel) context.Context
```

对应实现位于 `ctrl/ctrl.go` 的 `AppHandler` 接口定义。

### 变更原因与实现方式

`Event.Channel` 的处理是按 `Channel UUID` 串行投递到同一个事件线程中执行回调的。v1.4 中在 `ctrl/serve.go` 的 channel 事件处理逻辑里，会把回调返回的 `ctx` 重新赋值回当前闭包变量：

```go
ctx = handler.ChannelEvent(ctx, channelEvent)
```

这样对于同一个 channel 后续收到的事件回调，会继续拿到“上一次回调返回的 ctx”，从而能够拿到之前 `ctx` 里放入的数据。

### 迁移方法

1) 原有项目

把你的 Handler 方法签名改为带返回值，并在函数末尾 `return ctx` 即可：

```go
func (h *MyAppHandler) ChannelEvent(ctx context.Context, channel *ctrl.Channel) context.Context {
	// 原有逻辑保持不变
	return ctx
}
```

2) 需要在后续事件读取 START 阶段写入的自定义数据

在需要写入的时机用 `context.WithValue` 生成新 `ctx`，并把新 `ctx` 返回；后续事件再从 `ctx.Value` 取出即可：

```go
type ctxKey string

const userKey ctxKey = "user"

func (h *MyAppHandler) ChannelEvent(ctx context.Context, channel *ctrl.Channel) context.Context {
	switch channel.GetState() {
	case "START":
		ctx = context.WithValue(ctx, userKey, "alice")
		return ctx
	default:
		_ = ctx.Value(userKey)
		return ctx
	}
}
```

