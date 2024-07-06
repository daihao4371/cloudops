package models

import "github.com/golang-jwt/jwt/v5"

type UserLoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=25"` // 用户名格式 3-25位
	Password string `json:"password" validate:"required,min=3,max=25"` // 密码 3-25位
	Email    string `json:"email" validate:"email"`                    // 邮箱格式
	Gender   string `json:"gender" validate:"oneof=male female"`       // 必须在选择范围内                // 性别
}

// 根据token parse生成的对象
type UserCustomClaims struct {
	*User
	jwt.RegisteredClaims // 默认的注册字段
}

// 准备一个login接口的返回数据
type UserLoginResponse struct {
	*User
	Token     string `json:"token"` // 返回的token
	ExpiresAt int64  `json:"expiresAt"`
}
