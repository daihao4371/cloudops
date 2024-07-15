package models

import (
	"cloudops/src/common"
	"fmt"
	"gorm.io/gorm"
)

// 用户相关的数据库字段
type User struct {
	gorm.Model
	OrderNo  int     `json:"orderNo" gorm:"comment:用户id"`
	Username string  `json:"username" gorm:"type:varchar(100);uniqueIndex;comment:用户登录名"`
	Password string  `json:"password"  gorm:"comment:用户登录密码"`
	RealName string  `json:"realName" gorm:"type:varchar(100);uniqueIndex;comment:用户昵称 中文名"`
	Desc     string  `json:"desc" gorm:"comment:用户描述"`
	HomePath string  `json:"homePath" gorm:"comment:登陆后的默认首页"`
	Roles    []*Role `json:"roles" gorm:"many2many:user_roles;"`
	Enable   int     `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"`
}

/*
	func (u *User) TableName() string {
		return "users"
	}
*/
//CheckUserPassword 验证用户登录密码
func CheckUserPassword(ru *UserLoginRequest) (*User, error) {
	var dbUser User
	// 根据用户名查询用户信息，并预加载用户角色
	err := DB.Where("username = ?", ru.Username).Preload("Roles").First(&dbUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户名不存在")
		}
		return nil, fmt.Errorf("数据库错误:%w", err)
	}

	// 对比用户密码
	ok := common.BcryptCheck(ru.Password, dbUser.Password)
	if ok {

		return &dbUser, nil
	}
	return nil, fmt.Errorf("密码错误")
}

// GetUserByUserName 根据用户名获取用户信息
func GetUserByUserName(userName string) (*User, error) {
	var dbUser User
	// 根据用户名查询用户信息，并预加载用户角色及其菜单
	err := DB.Where("username = ?", userName).Preload("Roles.Menus").First(&dbUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户名不存在")
		}
		return nil, fmt.Errorf("数据库错误:%w", err)
	}
	return &dbUser, nil
}
