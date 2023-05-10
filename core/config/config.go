package config

import (
	"git.xswitch.cn/xswitch/xctrl/core/ctrl/store"
)

// fromCache .
func fromCache(key string) (string, error) {
	return store.NewConfig().Read(key)
}
