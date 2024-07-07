package models

import (
	"cloudops/src/config"
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
	)
}
