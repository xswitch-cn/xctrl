# NATS多租户实现 {#multi-tenant}

NATS支持多租户。

管理级租户名以`foo`为例，普通租户以`cherry`和`xyt`为例子。

在`foo`中接收来自`cherry`的消息，订阅以下主题：

```sh
*.cn.xswitch.ctrl
```

`cherry`租户向`cn.xswitch.ctrl`发消息，`foo`收则到后消息主题变为：

```sh
from-cherry.cn.xswitch.ctrl
```

同理，`foo`发送给`cherry`的消息，发送到以下主题：

```sh
to-cherry.cn.xswitch.node
```

`cherry`收到后消息主题变为`cn.xswitch.node`。

## 设置From和To Prefix

为了使消息更明确，可以设置明确的From和To Prefix，这样可以在NATS中更好的区分消息。

```go
ctrl.SetFromPrefix("from-")
ctrl.SetToPrefix("to-")
```

如果不设置，默认为空。要跟NATS服务端的相关配置一致。

## 订阅全部租户消息

订阅以下主题：

```go
*.cn.xswitch.ctrl
```

## ctrl.GetTenantId

`GetTenantId(subject string)`

该函数用于从`subject`中获取租户名。

如果`subject`中在`cn.xswitch.`之前有字符串，则认为是租户名，如：

```go
ctrl.GetTenantId("cn.xswitch.ctrl") == ""
ctrl.GetTenantId("cherry.cn.xswitch.ctrl") == "cherry"
ctrl.GetTenantId("from-cherry.cn.xswitch.ctrl") == "from-cherry"
ctrl.SetFromPrefix("from-")
ctrl.GetTenantId("from-cherry.cn.xswitch.ctrl") == "cherry"
```

注意，`SetFromPrefix`后可以使该函数去掉相关前缀。

## ChannelEvent

如果是`ChannelEvent`，在多租户情况可以使用`ChannelEvent.GetTenant()`获取租户名。

所有跟`ChannelEvent`（`ctrl.Channel`）相关的接口都无需明确使用租户名，系统收到消息后，在生成`ChannelEvent`之前会检查是该消息是否是在多租户方式来来自某一租户，并自动设置`ChannelEvent`的租户名。调用`ChannelEvent`相关接口（如`ChannelEvent.Answer()`）时，会自动加上租户名。

当然，也可以明确的设置租户名，如：

```go
tenant := channel.GetTenant()
option := channel.WithTentantAddress(tenant, "node-uuid-1")
channel.Answer(option)
```

### channel.WithTenantAddress

```go
WithTenantAddress(tenant string, address string) Option
```

设置租户名和地址，用于生成接口调用的`Option`。

示例：

```go
option := WithTenantAddress("cherry", "node-uuid-1")
```

### channel.GetTenant

获取`ChannelEvent`事件中的租户名，如果非多租户模式或消息是来自同一租户，则返回空字符串`""`。

## EnableNodeStatus

在多租户情况下，EnableNodeStatus会取到所有节点，并不会按租户区分。如果需要区分租户，则可以通过相应的回调在上层实现。

## cMan

在多租户状态下，类似`ctrl.InitCManService("to-cherry.cn.xswitch.cman.control")`中的Subject参数不再有用。在调用所有cMan相关接口时都需要明确传入租户名对应的主题。可以使用`ctrl.WithTenantAddress(tenant, address)`函数生成对应的Option。

如：

```go
	option := ctrl.WithTenantAddress("cherry", "cn.xswitch.cman.control")
	res, err := ctrl.CManService().GetConferenceList(context.Background(),
        &cman.GetConferenceListRequest{},
		ctrl.WithRequestTimeout(1*time.Second), option)
```

## 非ChannelEvent消息

非`ChannelEvent`消息将以`nats.Event`的形式传递给回调函数，可以通过`natsEvent.Topic()`获取原始的主题，并通过`ctrl.GetTenantID()`获取租户名，如：

```go
func (h *EventHandler) Event(req *ctrl.Request, natsEvent nats.Event) {
	tenant := ctrl.GetTenantID(natsEvent.Topic())
}
```
