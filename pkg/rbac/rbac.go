package rbac

import "github.com/casbin/casbin/v2/persist"

type RBAC struct {
	Adapter    persist.Adapter
	ConfigPath string
}

func NewRBAC(adapter persist.Adapter, ConfigPath string) RBAC {
	return RBAC{Adapter: adapter, ConfigPath: ConfigPath}
}
