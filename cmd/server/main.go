package main

import (
	"cloudops/src/common"
	"cloudops/src/config"
	"cloudops/src/models"
	"cloudops/src/web"
	"flag"
	"fmt"
	"go.uber.org/zap"
)

func main() {
	// 主配置文件 命令行解析
	var (
		configFile string
	)

	flag.StringVar(&configFile, "f", "./server.yml", "config file")
	flag.Parse()

	// 加载主配置文件
	sc, err := config.LoadServer(configFile)
	if err != nil {
		fmt.Printf("Loading main configuration file failed： %v\n", err)
		return
	}
	fmt.Printf("Load main configuration file success： %v\n", sc)
	// 初始化日志
	logger := common.NewZaplogger(sc.LogLevel, sc.LogPath, sc.ErrLogPath)
	defer logger.Sync()
	sc.Logger = logger

	// 初始化数据库
	err = models.InitDb(sc)
	if err != nil {
		logger.Error("InitDb failed", zap.String("err", err.Error()))
		return
	}
	logger.Info("InitDb success")

	/*	// 同步表结构
		err = models.MigrateTable()
		if err != nil {
			logger.Error("MigrateTable failed", zap.String("err", err.Error()))
			return
		}
		logger.Info("MigrateTable success")*/

	// 启动web服务
	err = web.StartGIn(sc)
}
