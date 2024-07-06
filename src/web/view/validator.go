package view

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhs "github.com/go-playground/validator/v10/translations/zh"
)

var (
	validate = validator.New()         //实例化验证器
	chinese  = zh.New()                // 获取中文翻译器
	uni      = ut.New(chinese)         // 设置成中文翻译器
	trans, _ = uni.GetTranslator("zh") // 获取翻译字典
	// 注册翻译器
	_ = zhs.RegisterDefaultTranslations(validate, trans)
)
