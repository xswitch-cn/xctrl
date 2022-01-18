package perm

import (
	"github.com/casbin/casbin/v2"
	//
	_ "github.com/go-sql-driver/mysql"
)

// Enforcer 权限管理器
var Enforcer *casbin.SyncedEnforcer


const rootUID = "c27a7757-f2a0-40b0-ab58-8eb5fe21b368"

// RootUID  超级用户UID
func RootUID() string {
	return rootUID
}


