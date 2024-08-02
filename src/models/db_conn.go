package models

import (
	"cloudops/src/config"
	"github.com/casbin/casbin/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB             *gorm.DB
	casbinEnforcer *casbin.Enforcer
)

// 初始化数据库连接
func InitDb(sc *config.ServerConfig) error {
	db, err := gorm.Open(
		mysql.Open(sc.MySqlC.DSN),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Info)},
		//&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)},
	)
	if err != nil {
		return err
	}
	DB = db
	return nil
}

// 手动注入SQL表结构
/*func MigrateTable() error {
	return DB.AutoMigrate(
		&User{},
		&Role{},
		&Menu{},
	)
}*/
