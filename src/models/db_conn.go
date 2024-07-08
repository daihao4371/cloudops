package models

import (
	"cloudops/src/common"
	"cloudops/src/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func InitDb(sc *config.ServerConfig) error {
	db, err := gorm.Open(
		mysql.Open(sc.MySqlC.DSN),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Info)},
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

func MockUserRegister(sc *config.ServerConfig) {
	u1 := User{
		Username: "admin",
		Password: "123456",
		RealName: "admin",
		Desc:     "",
		HomePath: "/system/account",
		Enable:   1,
		Roles: []*Role{
			{
				RoleName:  "超级管理员",
				RoleValue: "super",
			}, {
				RoleName:  "前端权限管理员",
				RoleValue: "fronAdmin",
			},
		},
	}
	u1.Password = common.BcryptHash(u1.Password)
	err := DB.Create(&u1).Error
	if err != nil {
		sc.Logger.Error("模拟注册用户失败", zap.String("错误", err.Error()))
	}
	sc.Logger.Info("模拟注册用户成功")
}
