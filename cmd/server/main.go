package main

import (
	"cloudops/src/common"
	"cloudops/src/config"
	"cloudops/src/web"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"log"
	"time"
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
		log.Fatalf("加载主配置文件失败： %v", err)
		return
	}
	logger := common.NewZaplogger(sc.LogLevel)
	defer logger.Sync()
	logger = common.NewZaplogger("DEBUG")
	logger.Info("failed to fetch URL",
		zap.String("url", "logger.Info"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

	logger.Debug("这是一条测试debug日志",
		zap.String("级别", sc.LogLevel),
	)

	fmt.Printf("主配置文件路径：%#v sc:%#v\n", configFile, sc)
	//log.Printf("主配置文件路径： %v  sc:%#v\n", configFile, sc)
	err = web.StartGIn(sc)
}
