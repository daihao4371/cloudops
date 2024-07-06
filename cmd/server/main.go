package main

import (
	"cloudops/src/common"
	"cloudops/src/config"
	"cloudops/src/web"
	"flag"
	"fmt"
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
		fmt.Printf("Loading main configuration file failed： %v", err)
		return
	}

	logger := common.NewZaplogger(sc.LogLevel, sc.LogPath, sc.ErrLogPath)
	defer logger.Sync()
	sc.Logger = logger

	/*	logger.Debug("这是一条测试debug日志",
			zap.String("级别", sc.LogLevel),
		)

		logger.Error("这是一条测试error日志")
	*/
	// 把logger设置到config中

	fmt.Printf("Main configuration file path：%#v sc:%#v\n", configFile, sc)
	//log.Printf("主配置文件路径： %v  sc:%#v\n", configFile, sc)
	err = web.StartGIn(sc)
}
