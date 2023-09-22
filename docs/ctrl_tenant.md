# NATS多租户实现 {#multi-tenant}

NATS支持多租户。

管理级租户名以`foo`为例，普通租户以`cherry`和`xyt`为例子。

在`foo`中接收来自`cherry`的消息，订阅以下主题：

```sh
cn.xswitch.ctrl
```

收则到消息后变为：

```sh
from-cherry.cn.xswitch.ctrl
```

同下，发送给`cherry`的消息，发送到以下主题：

```sh
to-cherry.cn.xswitch.cn
```

## 设置From和To Prefix

为了使消息更明确，可以设置明确的From和To Prefix，这样可以在NATS中更好的区分消息。

```go
ctrl.SetFromPrefix("from-")
ctrl.SetToPrefix("to-")
```

如果不设置，默认为空。要跟NATS服务端的相关配置一致。

## 订阅全部租户消息

```go
*.cn.xswitch.ctrl
```

## ChannelEvent

如果是`ChannelEvent`，在多租户情况可以使用`ChannelEvent.GetTenant()`获取租户名。

所有跟`ChannelEvent`（`ctrl.Channel`）相关的接口都需要明确使用租户名，如：

```go
tenant := channel.GetTenant()
option := channel.WithTentantAddress(tenant, "node-uuid-1")
channel.Answer(option)
```

### channel.WithTenantAddress

```go
WithTenantAddress(tenant string, address string) Option
```

示例：

```go
option := WithTenantAddress("cherry", "node-uuid-1")
```

### channel.GetTenant

获取事件中的租户名。

## cMan

在多租户状态下，`ctrl.InitCManService("to-cherry.cn.xswitch.cman.control")`中的参数不再有用。在调用所有cMan相关接口时都需要明确传入租户名。使用`ctrl.WithTenantAddress(tenant, address)`函数可以生成对应的Option。

如：

```go
	option := ctrl.WithTenantAddress("cherry", "cn.xswitch.cman.control")
	res, err := ctrl.CManService().GetConferenceList(context.Background(),
        &cman.GetConferenceListRequest{},
		ctrl.WithRequestTimeout(1*time.Second), option)
```
