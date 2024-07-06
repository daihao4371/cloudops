package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

// ServerConfig 定义服务器配置结构体

type ServeConfig struct {
	HttpAddr   string `yaml:"http_addr"`    // HTTP 服务器监听地址
	LogLevel   string `yaml:"log_level"`    // 日志级别
	LogPath    string `yaml:"log_path"`     // 日志文件路径
	ErrLogPath string `yaml:"err_log_path"` // 错误日志文件路径
}

// 根据IO read 读取配置文件后的字符串解析yaml
func LoadServer(filename string) (*ServeConfig, error) {
	cfg := &ServeConfig{}
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
