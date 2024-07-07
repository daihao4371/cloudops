package models

import (
	"gorm.io/gorm"
)

// 用户相关的数据库字段
type User struct {
	gorm.Model        // 模型基类
	UserId     int    `json:"UserId" gorm:"comment:用户id"`          // 用户id
	Username   string `json:"username" gorm:"index;comment:用户登录名"` // 建表索引
	Password   string `json:"password" gorm:"comment:用户登录密码"`
	RealName   string `json:"realName" gorm:"comment:用户昵称"`
	desc       string `json:"desc" gorm:"comment:用户描述"`
	HomePath   string `json:"homePath" gorm:"comment:登录后的默认首页"`
	//Role       []*Role `json:"-" gorm:"many2many:user_role;"` // 用户角色 关联表多对多
	//Phone      string  `json:"phone" gorm:"comment:手机号"`                        // 手机号
	//Email      string  `json:"email" gorm:"comment:邮箱"`                         // 邮箱
	Enable int `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"` // 用户是否启用	1正常 2冻结
}
