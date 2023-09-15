# ctrl 多实例说明

ctrl instance 是ctrl的多实例版本

Ctrl有关信息请参阅[ctrl.md](ctrl.md)

在绝大部分场景下和Ctrl的使用方法一样

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