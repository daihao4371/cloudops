package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	OrderNo   int     `json:"orderNo" gorm:"comment:排序"`
	RoleName  string  `json:"roleName" gorm:"type:varchar(100);uniqueIndex;comment:角色中文名称"`
	RoleValue string  `json:"roleValue"  gorm:"type:varchar(100);uniqueIndex;comment:角色值"`
	Remark    string  `json:"remark" gorm:"comment:用户描述"`
	HomePath  string  `json:"homePath" gorm:"comment:登陆后的默认首页"`
	Status    string  `json:"status" gorm:"default:1;comment:角色是否被冻结 1正常 2冻结"`
	Users     []*User `json:"users" gorm:"many2many:user_roles;"`
	Menus     []*Menu `json:"menus" gorm:"many2many:role_menus;"` // 多对多关系，角色对应的菜单
}
