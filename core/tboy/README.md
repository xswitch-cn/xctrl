# Test Boy

TBoy是一个测试框架，模拟FreeSWITCH的行为。

核心模拟一些基本的操作，可以直接使用，也可以根据需要进行扩展。参见`ecc/ttboy`的扩展用法。

考虑呼叫的各种情况，秒接，30秒接，不接，拒接，等等，都可以扩展tboy。我们吸收了多种情况下的tboy，在测试的时候选个合适的boy,给定适当的参数配置就可以测试多种情况，而不用自己拿着3个话机互相打。

## 点对点呼叫

```
node_uuid    = "simple-test.simple"
domain       = "test.test"
```

* 秒接

  ```
  boy := tboy.NewSimple(node_uuid, domain, tboy.OptionPeerAnswer(true))
  ```

* 30秒接

  ```
  boy := tboy.NewSimple(node_uuid, domain, tboy.OptionPeerAnswer(true), tboy.OptionPeerWait(30))
  ```

* 拒接

  ```
  boy := tboy.NewSimple(node_uuid, domain, tboy.OptionPeerReject(true))
  ```

  

## 队列呼叫

* 未分配坐席

  ```
  boy := tboy.NewACD(node_uuid, domain, tboy.OptionAcdAssign(false))
  ```

* 分配坐席未接

  ```
  boy := tboy.NewACD(node_uuid, domain, tboy.OptionAcdAssign(true), tboy.OptionPeerAnswer(false))
  ```

*  分配坐席接听

  ```
  boy := tboy.NewACD(node_uuid, domain, tboy.OptionAcdAssign(true), tboy.OptionPeerAnswer(true))
  ```

  

## 使用方式

1，预先启动一个nats服务器

2，根据呼叫测试情况，启用一个呼叫boy

```
node_uuid    = "simple-test.simple"
node_topic   = "simple-test"
domain       = "test.test"
boy := tboy.NewSimple(node_uuid, domain, tboy.OptionPeerAnswer(true))
err := ctrl.Init(boy, true, natsAddress)
	fmt.Println(err)

ctrl.Subscribe("cn.xswitch.node.test", boy.EventCallback, "node")
ctrl.Subscribe("cn.xswitch.node.test."+node_topic, boy.EventCallback, "node")
ctrl.Subscribe("cn.xswitch.node.test."+node_uuid, boy.EventCallback, "")
```

3，调用dial接口测试呼叫流程


nodejs调用单个方法


```
var rpc = {
	jsonrpc: "2.0",
	method: "XNode.Play",
	id: id,
	params: {
		ctrl_uuid: 'test-nodejs-controller',
		data: file
		uuid: uuid
	}
}
service = "cn.xswitch.node.test"
const msg = JSON.stringify(rpc);
console.log("sending " + msg);
nc.subscribe('cn.xswitch.ctrl.test-nodejs-controller',function (msg, reply, subject, sid) {...});
nc.request(service, msg, { max: 1, timeout: 5000 }, (msg) => {...});
```

