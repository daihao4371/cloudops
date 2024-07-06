package models

// 用户相关的数据库字段
type User struct {
	//Id       int    `json:"id"`
	Username string `json:"username" gorm:"index;comment:用户登录名"`
	Password string `json:"-" gorm:"comment:用户登录密码"`
}
