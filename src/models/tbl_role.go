package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	OrderNo   int64   `json:"order_no" gorm:"comment:排序号"`
	RoleName  string  `json:"roleName" gorm:"type:varchar(100);uniqueIndex;comment:角色中文名称"`
	RoleValue string  `json:"roleValue" gorm:"type:varchar(100);uniqueIndex;comment:角色值"`
	Remark    string  `json:"remark" gorm:"type:longtext;comment:角色描述"`
	HomePath  string  `json:"homePath" gorm:"type:longtext;comment:登陆后的默认首页"`
	Status    string  `json:"status" gorm:"type:varchar(191);default:1;comment:角色是否被冻结 1正常 2冻结"`
	Users     []*User `json:"users" gorm:"many2many:user_roles"`
}

func (r *Role) TableName() string {
	return "roles"
}
