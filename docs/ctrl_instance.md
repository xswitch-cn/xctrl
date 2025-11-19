# ctrl 多实例说明 {#multi-instance}

ctrl 可以支持多实例，以下列举了一些常用的方法。

**注意**： 一般情况下不需要用到多实例，如果可以用NATS多租户方式解决，就不要用多实例。

## 使用说明

```go
func NewCtrlInstance(trace bool, addrs string) (*Ctrl, error)
```
- 默认会返回 ctrl实例 指针和error，需要自行判断error

## 代码示例

```go
instance, err := ctrl.NewCtrlInstance(true, natsURL)
if err != nil {
    log.Error("Ctrl isn't found")
}
instance.EnableEvent(new(CtrlInstanceEvent1), "test.test", "")

```

## SetMaxChannelLifeTime

```go
func(c *Ctrl) SetMaxChannelLifeTime(time uint)
```

设置Channel的最长保留事件小时数，超时就会销毁，默认为4小时

* `time`：为小时数

## SetMaxChannelLifeTimeMinute

```go
func(c *Ctrl) SetMaxChannelLifeTimeMinute(time uint)
```

设置Channel的最长保留事件分钟数，超时就会销毁，默认为4小时

* `time`：分钟数

## EnableApp

```go
fun(c *Ctrl)c EnableApp(h *AppHandler, subject string, queue string) error
```

订阅一个`Topic`，是一个全能的订阅方式，包括接收Node的事件、返回结果等。如果`subject`为`cn.xswitch.ctrl`，则它会订阅两个主题：

- `cn.xswitch.ctrl`：用于接收Node的事件，以队列方式订阅，主要用于多个ctrl负载分担方式获取事件消息。
- `cn.xswitch.ctrl.$ctrl_uuid`：用于接收Node的事件，非队列方式订阅，仅接受本ctrl自己的消息。其中`$ctrl_uuid`为自动生成的UUID，是从`ctrl.UUID()`获取的。

参数含义如下：

- `h`是一个`ctrl.AppHandler`类型的结构，必须实现它定义的几个函数，下面会有详细描述。
- `subject`和`queue`是订阅的NATS主题和队列，如果`queue`为空，则不使用队列方式订阅，详见NATS对Request/Reply模式的相关说明。

Handler是一个`interface`，必须实现如下几个函数（可以是空函数）。

```go
type Handler interface {
	ChannelEvent(ctx context.Context, channel *Channel)
	Event(msg *Message, natsEvent nats.Event)
}
```

对于`Event.Channel`事件，回调函数里它将以当前的channel的`uuid`为topic和queue启用一个`bus`消息总线进行订阅处理，一方面避免NATS回调端的阻塞，另一方面，使channel在bus中成为一个串行的订阅。因而，对于同一个Channel UUID来说，回调是串行的，保证channel的`START`、`RING`、`ANSWER`、`DESTROY`等事件处理的有序性。

对于其它事件，它将使用新的Go Routine进行回调，因而，无法保证顺序。

- `ctx`：对于`ChannelEvent`，`ctx`是个`context.Context`，可以从里面取到最原始的Channel信息（第一个Channel事件，如`START`或`READY`）。

    ```go
    var key ctrl.ContextKey = "channel"
    channel := ctx.Value(key).(*ctrl.Channel)
    ```

- `channel`：对于次收到的`Event.Channel`事件，都将转换成一个新的`ctrl.Channel`类型的结构体。

可以通过`channel.GetNatsEvent()`获取原始的NATS事件。

对于非Channel的事件，都在`Event`函数中回调。

- `msg`：JSON RPC格式的消息，可能是个请求，也可能是个响应。
    - 如果消息中有`msg.Result`，说明是个正常响应。
    - 如果消息中有`msg.Error`，说明是个错误响应。
    - 否则是个请求。
        - 如果有`msg.ID`，说明是个有ID的请求，应该返回一个响应。
        - 否则是个事件。
- `natsEvent`：原始的NATS消息。

一般来说一个应用程序仅调用一次`EnableApp`，对多个`EnableApp`的调用，结果是未知的。

### 示例

```go
instance, err := ctrl.NewCtrlInstance(true,natsUrl)

subject := "cn.xswitch.ctrl"
type AppExample struct {}

func (h *AppExample) Event(msg *ctrl.Message, natsEvent nats.Event) {}

func (a *AppExample) ChannelEvent(ctx context.Context, c *ctrl.Channel) {}

instance.EnableApp(new(AppExample),subject,"")
```

### EnableEvent

```go
func(c *Ctrl) EnableEvent(h *EventHandler, subject string, queue string) error
```

订阅事件对应的Subject，如`cn.xswitch.ctrl.cdr`。目前，除`cn.xswitch.ctrl.cdr`是在NATS中串行回调外，其它均为在新的Go Routine中回调。

如果一个应用程序中即调用`EnableApp`和`EnableEvent`，则两者的`subject`不要重复，否则会有不可预知的结果。

### 示例

```go
instance, err := ctrl.NewCtrlInstance(true,natsUrl)

subject := "cn.xswitch.ctrl"
type EventExample struct {}

func (h *EventExample) Event(req *ctrl.Request, natsEvent nats.Event) {}

instance.EnableEvent(new(EventExample), subject, "")
```

### EnableRequest

```go
func(c *Ctrl) EnableRequest(h *RequestHandler, subject string, queue string) error
```

订阅Request请求消息。主要用于处理FreeSWITCH的XML或JSON数据配置请求，如`dialplan`、`directory`、`config`等。这种订阅总是异步处理的。

当然也可以处理通用的请求。

如果一个应用程序中即调用`EnableApp`和`EnableRequest`，则两者的`subject`不要重复，否则会有不可预知的结果。

### 示例

```go
instance, err := ctrl.NewCtrlInstance(true,natsUrl)

subject := "cn.xswitch.ctrl"

type RequestExample struct {}

func (r RequestExample) Request(req *ctrl.Request, natsEvent nats.Event)  {}

instance.EnableRequest(new(RequestExample),subject,"")

```

### EnableNodeStatus

```go
func(c *Ctrl) EnableNodeStatus(subject string) error
```

如果`subject`为空，则使用默认的`cn.xswitch.status.node`。

**注意**：在多ctrl的场景中，由于默认的订阅主题`cn.xswitch.ctrl`是通过队列方式订阅的，多个ctrl无法同时接收到节点状态，因此，需要使用独立的`EnableNodeStatus`订阅。

### 示例

```go
instance, err := ctrl.NewCtrlInstance(true, natsUrl)

subject := "cn.xswitch.status.node"
instance.EnableNodeStatus(subject)
```

### OnEvicted

```go
func(c *Ctrl) OnEvicted(f func(string, interface{}))
```

设置节点过期回调函数，如果Node节点过期，将会调用此回调函数。

目前内置定时器固定`10`秒检查一次，因此，最长可能在`Expiry`过期时间`10`秒后才能触发。

### 示例

```go
instance, err := ctrl.NewCtrlInstance(true,natsUrl)

instance.OnEvicted(func(s string, i interface{}){
    log.Printf("Node %s has expired", s)
})
```

### GetNATSConn

```go
func(c *Ctrl) GetNATSConn() *natsio.Conn
```

返回Ctrl内部的nats Connection的对象，用于修改内部默认方法。 具体请参考：

* [NATS Reconnect](https://docs.nats.io/using-nats/developer/connecting/reconnect)
* [Go Demo](../example/nats-conn/demo.go)

### 示例

```go
instance, err := ctrl.NewCtrlInstance(true,natsUrl)

instance.GetNATSConn()
```
### Subscribe

```go
func Subscribe(subject string, cb nats.EventCallback, queue string) (nats.Subscriber, error)
```

调用底层的NATS发起一个订阅。所有回调在同一个NATS Go Routine中回调。需要避免阻塞。

### Context

Ctrl中的Context使用了标准的Go Context包，目前没有太大用处，大部分可以直接传入`context.Background()`或`context.TODO()`。

### queueBufferSize

在订阅事件的时候会使用这个变量大小进行channel的初始化，1024容量足够事件使用，太小会导致程序阻塞卡顿，影响运行效率。

### protobuf 扩展示例

```go
req := &xctrl.XNativeJSRequest{
	CtrlUuid: CtrlUUID,
	Data: &xctrl.XNativeJSRequestData{
		Command: "sofia.status",
		Data: *ctrl.ToRawMessage(map[string]string{
			"profile": profile_name,
		}),
	},
}
instance, err := ctrl.NewCtrlInstance(true,natsUrl)
response, err := instance.Service().NativeJSAPI(context.Background(), req, ctrl.WithAddress(""))
```
