package models

func CasbinCheckPermission(roleName, path, method string) (bool, error) {
	return casbinEnforcer.Enforce(roleName, path, method)
}

func CasbinAddOnePolicy(sub, obj, act string) (bool, error) {
	return casbinEnforcer.AddPolicy(sub, obj, act)
}
