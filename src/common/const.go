package common

/*
// GIN_CTX_CONFIG_LOGGER 是存储在 Gin 上下文中的日志配置键
const GIN_CTX_CONFIG_LOGGER = "gin_logger"

// GIN_CTX_CONFIG_CONFIG 是存储在 Gin 上下文中的应用配置键
const GIN_CTX_CONFIG_CONFIG = "gin_config"

// GIN_CTX_JWT_CLAIM 是存储在 Gin 上下文中的 JWT 声明键
const GIN_CTX_JWT_CLAIM = "jwt_claim"
*/
const (
	GIN_CTX_CONFIG_LOGGER = "gin_logger"
	GIN_CTX_CONFIG_CONFIG = "gin_config"
	GIN_CTX_JWT_CLAIM     = "jwt_claim"
	GIN_CTX_JWT_USER_NAME = "jwt_user_name"

	COMMON_STATUS_ENABLE  = "1"
	COMMON_STATUS_DISABLE = "0"
	GIN_CTX_STREE_CACHE   = "stree_cache"
)

var (
	COMMON_SHOW_MAP = map[string]bool{
		"1": true,
		"0": false,
	}
)
