package models

type UserLoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=25"`         // 用户名格式 3-25位
	Password string `json:"password" validate:"required,min=3,max=25"`         // 密码 3-25位
	Email    string `json:"email" validate:"email"`                            // 邮箱格式
	Gender   string `json:"gender" validate:"oneof=male female prefer_not_to"` // 必须在选择范围内                // 性别
}
