package models

import (
	"cloudops/src/common"
	"fmt"
	"gorm.io/gorm"
)

// 用户相关的数据库字段
type User struct {
	gorm.Model         // 模型基类
	UserId     int     `json:"UserId" gorm:"comment:用户id"`                                  // 用户id
	Username   string  `json:"username" gorm:"type:varchar(100);uniqueIndex;comment:用户登录名"` // 用户登录名 uniqueIndex 代表唯一索引
	Password   string  `json:"password" gorm:"comment:用户登录密码"`
	RealName   string  `json:"realName" gorm:"comment:用户昵称"`
	desc       string  `json:"desc" gorm:"comment:用户描述"`
	HomePath   string  `json:"homePath" gorm:"comment:登录后的默认首页"`
	Role       []*Role `json:"roles" gorm:"many2many:user_roles;"` // 用户角色 关联表多对多
	//Phone      string  `json:"phone" gorm:"comment:手机号"`                        // 手机号
	//Email      string  `json:"email" gorm:"comment:邮箱"`                         // 邮箱
	Enable int `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"` // 用户是否启用	1正常 2冻结
}

func (u *User) TableName() string {
	return "users"
}

func CheckUserPassword(ru *UserLoginRequest) (*User, error) {
	// 查询用户
	dbUser := User{
		Username: ru.Username,
		Password: ru.Password,
	}
	err := DB.First(&dbUser).Error
	if err != nil {
		return nil, err
	}

	//  验证密码
	ok := common.BcryptCheck(ru.Password, dbUser.Password)
	if ok {
		return &dbUser, nil
	}
	return nil, fmt.Errorf("密码错误")
}
