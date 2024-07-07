package models

import (
	"cloudops/src/common"
	"cloudops/src/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func InitDb(sc *config.ServeConfig) error {
	db, err := gorm.Open(
		mysql.Open(sc.MySqlC.DSN),
		&gorm.Config{},
	)
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func MigrateTable() error {
	return DB.AutoMigrate(
		&User{},
		&Role{},
	)
}

func MockUserRegister(sc *config.ServeConfig) {
	u1 := User{
		Model:    gorm.Model{},
		UserId:   0,
		Username: "admin",
		Password: "123456",
		RealName: "admin",
		desc:     "超级管理员",
		HomePath: "/system/account",
		Enable:   1,
		Role: []*Role{
			{
				RoleName:  "超级管理员",
				RoleValue: "super",
			}, {
				RoleName:  "前端管理员",
				RoleValue: "fronAdmin",
			},
		},
	}
	u1.Password = common.BcrypaHash(u1.Password)
	err := DB.Create(&u1).Error
	if err != nil {
		sc.Logger.Error("模拟注册用户失败", zap.Any("err", err))
	}
	sc.Logger.Info("模拟注册用户成功", zap.Any("user", u1))
}
