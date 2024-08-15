package cache

import (
	"cloudops/src/config"
	"sync"
)

// 关于监控的缓存
type StreeCache struct {
	// 主配置文件map key是Prometheus ip value 是主配置文件
	StreeNodeCacahe sync.Map
	Sc              *config.ServerConfig
	sync.RWMutex
	AlertLock        sync.RWMutex
	MainLock         sync.RWMutex
	AlertManagerLock sync.RWMutex
	RecordLock       sync.RWMutex
}

func NewStreeCache(sc *config.ServerConfig) *StreeCache {
	mc := &StreeCache{
		StreeNodeCacahe: sync.Map{},
		Sc:              sc,
		RWMutex:         sync.RWMutex{},
	}
	return mc
}
