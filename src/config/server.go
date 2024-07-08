package config

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"os"
	"time"
)

// ServerConfig 定义服务器配置结构体
type ServerConfig struct {
	HttpAddr   string        `yaml:"http_addr"`    // HTTP 服务器监听地址
	MySqlC     *mysql.Config `yaml:"mysql"`        // MySQL 配置
	LogLevel   string        `yaml:"log_level"`    // 日志级别
	LogPath    string        `yaml:"log_path"`     // 日志文件路径
	ErrLogPath string        `yaml:"err_log_path"` // 错误日志文件路径
	JWTC       *JWT          `yaml:"jwt"`          // JWT配置
	Logger     *zap.Logger   `yaml:"-"`            // 不是配置文件中的字段，而是共用配置文件这个结构体，在生成完logger后设置的
}

type JWT struct {
	SingingKey      string        `yaml:"singing_key"`     // JWT签名 密码加盐
	ExpireTime      string        `yaml:"expire_time"`     // 过期时间
	ExpiresDuration time.Duration `yaml:"-"`               // 解析配置文件你的时候set的
	BufferTime      string        `yaml:"buffer_time"`     // 临期时间
	BufferDuration  time.Duration `yaml:"buffer_duration"` // 临期时间
	Issuers         string        `yaml:"issuers"`         // JWT签发者
}

// 根据IO read 读取配置文件后的字符串解析yaml
func LoadServer(filename string) (*ServerConfig, error) {
	cfg := &ServerConfig{}
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		return nil, err
	}

	// 解析JWTC的过期时间
	exd, err := time.ParseDuration(cfg.JWTC.ExpireTime)
	if err != nil {
		return nil, err
	}

	bud, err := time.ParseDuration(cfg.JWTC.ExpireTime)
	if err != nil {
		return nil, err
	}

	cfg.JWTC.ExpiresDuration = exd
	cfg.JWTC.BufferDuration = bud
	return cfg, nil
}
