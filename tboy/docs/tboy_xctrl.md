# TBoy XCtrl 使用示例

本文介绍如何使用 TBoy XCtrl 模拟 FreeSWITCH 通话，并与 XCtrl 控制器协同工作。
该示例展示了如何同时：

* 自动生成呼叫事件；
* 让多个 XCtrl 控制实例接收并响应；
* 模拟完整的通话过程（呼入 → 接听 → 挂断）。

---

## 环境准备

启动 NATS 服务

设置环境变量（可选）：

```dotenv
 NATS_URL="nats://127.0.0.1:4222"
```
---

## 示例代码结构

完整示例：

```go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
	"os"
	"os/signal"
	"syscall"

	"git.xswitch.cn/xswitch/proto/go/proto/xctrl"
	"git.xswitch.cn/xswitch/xctrl/ctrl"
	"git.xswitch.cn/xswitch/xctrl/ctrl/nats"
	"git.xswitch.cn/xswitch/xctrl/tboy"
	log "github.com/sirupsen/logrus"
)

const subject = "cn.xswitch.ctrl"

func main() {
	// 捕获退出信号
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// 读取 NATS 地址
	natsUrl := os.Getenv("NATS_URL")
	if natsUrl == "" {
		natsUrl = "nats://127.0.0.1:4222"
	}

	// 创建自动应答模拟器，每5秒生成一次新呼叫
	fakeAnswerChannel := tboy.NewFakeAnswerChannel(5*time.Second, natsUrl)
	fakeAnswerChannel.Start()

	// 启动测试控制器
	xTest := &XctrlTest{}
	go func() {
		for {
			// 最多保持 10 个并发测试控制器
			if xTest.GetCount() > 10 {
				time.Sleep(5 * time.Second)
				return
			}

			// 创建一个新的 XCtrl 实例
			ctrl1, err := ctrl.NewCtrlInstance(false, natsUrl)
			if err != nil {
				log.Error(err)
				return
			}

			// 将控制器 UUID 注册到 TBoy 模拟器
			fakeAnswerChannel.AddTopic(ctrl1.UUID())

			// 注册事件回调
			ctrl1.EnableApp(xTest, subject, "")

			xTest.AddTest()
		}
	}()

	<-shutdown
	fmt.Println("shutting down")
}
```

---

## 控制器逻辑说明

流程如下：

1. 收到 TBoy 生成的呼叫；
2. 模拟接听（Answer）；
3. 等待 10 秒；
4. 模拟挂机；
5. 从计数中移除该通话。

---

### 事件消息回调

```go
func (x *XctrlTest) Event(msg *ctrl.Message, natsEvent nats.Event) {
	log.Println("event:", msg)
}
```

该回调会打印所有从 TBoy 发送过来的事件。

---

### 控制器计数

为了防止无限创建，使用计数器控制并发上限：

```go
func (x *XctrlTest) AddTest() { ... }
func (x *XctrlTest) DelTest() { ... }
func (x *XctrlTest) GetCount() int { ... }
```

当前示例最多并行 10 个呼叫。

---

---

## 运行效果总结

| 功能     | 说明                |
| ------ | ----------------- |
| 自动生成呼叫 | 每隔 5 秒创建新通话       |
| 并发控制   | 最多并行 10 个控制器      |
| 自动应答   | 模拟接听后挂机           |
| 自动生成话单 | 每个通话结束自动发送CDR     |
| 优雅退出   | 支持 Ctrl+C 停止并清理资源 |

---

## 扩展应用场景

| 场景           | 示例                                  |
| ------------ | ----------------------------------- |
| 性能压测     | 调低 timeout（如 1 秒），提升呼叫速率            |
| 多控制节点测试  | 在 `AddTopic()` 注册多个 controller UUID |
| 长时间稳定性验证 | 运行数小时，观察 NATS 消息和内存变化               |
| 自定义通话行为  | 修改 `ChannelEvent()` 实现自定义挂机逻辑       |
| CI 集成测试  | 可嵌入自动化测试流程验证 XCtrl SDK 功能           |

---
