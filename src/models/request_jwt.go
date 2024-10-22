package models

import (
	"cloudops/src/common"
	"cloudops/src/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"time"
)

func TokenNext(dbUser *User, c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	// 生成JWTtoken
	token, err := GenJwtToken(dbUser, sc)
	if err != nil {
		sc.Logger.Error("Failed to generate JWT", zap.Error(err))
		common.FailWithMessage("Failed to generate jwt", c)
		return
	}
	// 构造返回结构体
	userResp := UserLoginResponse{
		User:      dbUser,
		Token:     token,
		ExpiresAt: 0,
	}
	common.OkWithDetailed(userResp, "User login successful", c)
}

// 生成JWT
func GenJwtToken(dbUser *User, sc *config.ServerConfig) (string, error) {
	// new claim对象
	c := UserCustomClaims{
		User: dbUser,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: sc.JWTC.Issuers, // 签发人
			// 默认5分钟过期：第一次生成的时候，过期时间戳，当前时间往后推 过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(sc.JWTC.ExpiresDuration)), // 过期时间
		},
	}
	// 使用制定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的签名密钥生成token字符串
	return token.SignedString([]byte(sc.JWTC.SingingKey))
}

// 解析token
func ParseToken(jwtLongToken string, sc *config.ServerConfig) (*UserCustomClaims, error) {
	// 解析token
	tokenClaims, err := jwt.ParseWithClaims(
		jwtLongToken, &UserCustomClaims{},
		func(token *jwt.Token) (i interface{}, e error) {
			return []byte(sc.JWTC.SingingKey), nil
		})
	if err != nil {
		sc.Logger.Error("Failed to parse token", zap.Error(err))
		return nil, err
	}
	// 判断是否有效
	if claims, ok := tokenClaims.Claims.(*UserCustomClaims); ok && tokenClaims.Valid {
		//sc.Logger.Error("Failed to parse token", zap.Any("过期时间", claims.RegisteredClaims.ExpiresAt))

		return claims, nil
	}
	return nil, err
}
