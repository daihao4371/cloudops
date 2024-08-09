package models

import (
	"cloudops/src/config"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
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

		// gorm 打印sql日志 开发模式使用
		//&gorm.Config{Logger: logger.Default.LogMode(logger.Info)},

		// gorm 打印sql日志 生产模式使用
		&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)},
	)
	if err != nil {
		return err
	}
	DB = db
	return nil
}

// 初始化casbin权限管理
func InitCasbin(sc *config.ServerConfig) error {
	a, err := gormadapter.NewAdapterByDB(DB)
	if err != nil {
		sc.Logger.Error("初始化casbin数据库错误", zap.Error(err))
		return err
	}
	// 初始化模型
	modelText := `
	[request_definition]
	r = sub, obj, act
	
	[policy_definition]
	p = sub, obj, act
	
	[role_definition]
	g = _, _
	
	[policy_effect]
	e = some(where (p.eft == allow))
	
	[matchers]
	m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
    `
	m, err := model.NewModelFromString(modelText)
	if err != nil {
		sc.Logger.Error("casbin字符串加载模型失败!", zap.Error(err))
		return err
	}

	casbinEnforcer, err = casbin.NewEnforcer(m, a)
	if err != nil {
		sc.Logger.Error("casbin创建Enforcer失败!", zap.Error(err))
		return err
	}
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
