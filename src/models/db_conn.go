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
		//&gorm.Config{Logger: logger.Default.LogMode(logger.Info)},
		&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)},
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
		&Menu{},
	)
}

func MockUserRegister(sc *config.ServerConfig) {
	menus := []*Menu{
		{
			Name:      "System",
			Title:     "系统管理",
			Icon:      "ion:settings-outline",
			Type:      "0",
			Show:      "1",
			OrderNo:   10,
			Component: "LAYOUT",
			Redirect:  "/system/account",
			Path:      "/system",
		},
		{
			Name:       "MenuManagement",
			Title:      "菜单管理",
			Icon:       "ant-design:account-book-filled",
			Type:       "1",
			Show:       "1",
			OrderNo:    11,
			Component:  "/demo/system/menu/index",
			ParentMenu: "1",
			Path:       "menu",
		},
		{
			Name:       "AccountManagement",
			Title:      "用户管理",
			Icon:       "ant-design:account-book-twotone",
			Type:       "1",
			Show:       "1",
			OrderNo:    12,
			Component:  "/demo/system/account/index",
			ParentMenu: "1",
			Path:       "account",
		},
		{
			Name:       "RoleManagement",
			Title:      "角色管理",
			Icon:       "ion:layers-outline",
			Type:       "1",
			Show:       "1",
			OrderNo:    13,
			Component:  "/demo/system/role/index",
			ParentMenu: "1",
			Path:       "role",
		},
		{
			Name:       "ChangePassword",
			Title:      "修改密码",
			Icon:       "ion:layers-outline",
			Type:       "1",
			Show:       "1",
			OrderNo:    14,
			Component:  "/demo/system/password/index",
			ParentMenu: "1",
			Path:       "changePassword",
		},
		{
			Name:       "Permission",
			Title:      "权限管理",
			Icon:       "ion:layers-outline",
			Type:       "0",
			Show:       "1",
			OrderNo:    14,
			Component:  "LAYOUT",
			ParentMenu: "1",
			Path:       "/permission",
			Redirect:   "/permission/front/page",
		},
		{
			Name:       "PermissionFrontDemo",
			Title:      "前端权限管理",
			Icon:       "ion:layers-outline",
			Type:       "1",
			Show:       "1",
			OrderNo:    15,
			Component:  "/demo/permission/index",
			ParentMenu: "6",
			Path:       "/front",
		},
	}

	/*	for _, menu := range menus {
		err := DB.Create(&menu).Error
		if err != nil {
			sc.Logger.Error("模拟注册菜单失败", zap.String("错误", err.Error()))
		}
		sc.Logger.Info("模拟注册菜单成功")
	}*/

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
				Menus:     menus,
			}, {
				RoleName:  "前端权限管理员",
				RoleValue: "fronAdmin",
				//Menus:     menus,
			},
		},
	}
	u1.Password = common.BcryptHash(u1.Password)
	err := DB.Create(&u1).Error
	if err != nil {
		sc.Logger.Error("模拟注册用户失败", zap.String("错误", err.Error()))
	}
	//DB.Save(&u1)
	sc.Logger.Info("模拟注册用户成功")
	sc.Logger.Info("模拟注册菜单成功")
}
