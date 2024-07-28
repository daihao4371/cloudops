package models

import "time"

type Model struct {
	ID          uint `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatedTime string `json:"createdTime" gorm:"-"`
	UpdatedTime string `json:"createdTime" gorm:"-"`
}

type EchartOneItem struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func commonGetUserNamesByUsers(users []*User) []string {
	userNames := []string{}
	for _, user := range users {
		user := user
		userNames = append(userNames, user.Username)
	}
	return userNames
}
