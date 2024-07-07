package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	OrderNo   int     `json:"order_no" gorm:"comment:排序号"`
	RoleName  string  `json:"roleName" gorm:"type:varchar(100);uniqueIndex;comment:角色名称"`
	RoleValue string  `json:"roleValue" gorm:"type:varchar(100);uniqueIndex;comment:角色值"`
	Reamrk    string  `json:"reamrk" gorm:"cocomment:用户描述"`
	HomePath  string  `json:"HomePath" gorm:"comment:登录后的默认首页"`
	Status    int     `json:"status" gorm:"comment:状态 0正常 1停用"`
	Users     []*User `gorm:"many2many:user_roles;"` // 多对多
}

func (r *Role) TableName() string {
	return "roles"
}
