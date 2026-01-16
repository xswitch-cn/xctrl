package store

import (
	"testing"
	"time"

	"github.com/xswitch-cn/proto/xctrl/store"
)

func init() {
	Init()
}

func TestStore(t *testing.T) {

}

func TestExpiry(t *testing.T) {
	Write("test", "test", store.WriteTTL(time.Second*30))
}

func TestDelete(t *testing.T) {
	// 路由
	NewRoute("dev.xswitch.cn").Delete()
	// 分机
	NewRegister("dev.xswitch.cn", "6200").Delete()
	// 模板
	NewTemplate("dialplan", "extension.conf", "default").Delete()
	// 全局配置
	// NewConfig().Delete(config.AuthExpiry)
}
