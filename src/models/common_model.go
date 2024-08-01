package models

import "time"

// 通用模型
type Model struct {
	ID          uint `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatedTime string `json:"createdTime" gorm:"-"`
	UpdatedTime string `json:"createdTime" gorm:"-"`
}

// ECharts图表中的一个数据项的结构体
type EchartOneItem struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// 根据用户列表返回用户名列表
func commonGetUserNamesByUsers(users []*User) []string {
	userNames := []string{}
	for _, user := range users {
		user := user
		userNames = append(userNames, user.Username)
	}
	return userNames
}
