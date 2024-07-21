package view

import (
	"cloudops/src/common"
	"cloudops/src/config"
	"cloudops/src/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func UserLogin(c *gin.Context) {
	// 校验用户名密码
	var user models.UserLoginRequest
	err := c.ShouldBindJSON(&user)

	// 判断JSON解析是否正确
	if err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	// 校验字段是否必填，范围是否正确
	err = validate.Struct(user)
	if err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			common.ReqBadFailWithWithDetailed(
				gin.H{
					"翻译前": err.Error(),
					"翻译后": errors.Translate(trans),
				},
				"Request error",
				c,
			)
			return
		}

		common.ReqBadFailWithMessage(err.Error(), c)
		return
	}

	dbUser, err := models.CheckUserPassword(&user)
	if err != nil {
		sc.Logger.Error("登录失败，用户名不存在或密码错误", zap.Error(err))
		common.ReqBadFailWithMessage(fmt.Sprintf("用户名不存在或密码错误:%v", err.Error()), c)
		return
	}
	// 生成token，并返回给前端
	models.TokenNext(dbUser, c)
}

// 登录以后获取用户信息
func getUserInfoAfterLogin(c *gin.Context) {

	/*	userName := c.MustGet(common.GIN_CTX_JWT_USER_NAME).(string)
		sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
		dbUser, err := models.GetUserByUserName(userName)
		if err != nil {
			sc.Logger.Error("通过token解析到的username去数据库找User失败", zap.Error(err))
			common.ReqBadFailWithMessage(fmt.Sprintf("通过token解析到的username去数据库找User失败:%v", err.Error()), c)
			return
		}
		common.OkWithDetailed(dbUser, "获取用户信息成功", c)*/

	userName, err := getContextUserName(c)
	if err != nil {
		handleError(c, err)
		return
	}
	sc, err := getContextServerConfig(c)
	if err != nil {
		handleError(c, err)
		return
	}
	dbUser, err := models.GetUserByUserName(userName)
	if err != nil {
		sc.Logger.Error("获取用户信息失败", zap.Error(err))
		common.FailWithMessage(fmt.Sprintf("获取用户信息失败:%v", err.Error()), c)
		return
	}
	common.OkWithDetailed(dbUser, "获取用户信息成功", c)
}

func getPermCode(c *gin.Context) {
	common.OkWithDetailed("123", "获取权限码成功", c)
}

func getContextUserName(c *gin.Context) (string, error) {
	userNameI, exists := c.Get(common.GIN_CTX_JWT_USER_NAME)
	if !exists {
		return "", fmt.Errorf("context中不存在用户名")
	}
	userName, ok := userNameI.(string)
	if !ok {
		return "", fmt.Errorf("context中的用户名类型不正确")
	}
	return userName, nil
}

func getContextServerConfig(c *gin.Context) (*config.ServerConfig, error) {
	scI, exists := c.Get(common.GIN_CTX_CONFIG_CONFIG)
	if !exists {
		return nil, fmt.Errorf("context中不存在服务器配置")
	}
	sc, ok := scI.(*config.ServerConfig)
	if !ok {
		return nil, fmt.Errorf("context中服务器配置类型错误")
	}
	return sc, nil
}

func handleError(c *gin.Context, err error) {
	common.ReqBadFailWithMessage(err.Error(), c)
}

func createAccount(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	// 校验一下 menu字段
	var reqUser models.User
	err := c.ShouldBindJSON(&reqUser)
	if err != nil {
		sc.Logger.Error("解析新增用户请求失败", zap.Any("用户", reqUser), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 在这里校验字段，是否必填，范围是否正确
	err = validate.Struct(reqUser)
	if err != nil {
		// 这里为什么要判断错误是否是 ValidationErrors
		if errors, ok := err.(validator.ValidationErrors); ok {
			common.ReqBadFailWithWithDetailed(
				gin.H{
					"翻译前": err.Error(),
					"翻译后": errors.Translate(trans),
				},
				"请求出错",
				c,
			)
			return
		}
		common.ReqBadFailWithMessage(err.Error(), c)
		return
	}

	sc.Logger.Info("创建用户请求字段打印", zap.Any("用户", reqUser))

	// 根据 rolesFront 去db中查询 role ，把role给他关联一下
	reqUser.Roles = make([]*models.Role, 0)
	for _, roleValue := range reqUser.RolesFront {
		dbRole, err := models.GetRoleByRoleValue(roleValue)
		if err != nil {
			sc.Logger.Error("根据RolesFront去db中找角色失败", zap.Any("用户", reqUser), zap.Error(err))
			common.FailWithMessage(err.Error(), c)
			return
		}
		reqUser.Roles = append(reqUser.Roles, dbRole)
	}

	// 直接存储加密后的密码
	err = reqUser.CreateOne()
	if err != nil {
		sc.Logger.Error("创建用户错误", zap.Any("用户", reqUser), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	common.OkWithMessage("创建成功", c)
}

func getAccountList(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	// 数据库中拿到所有的menu列表
	users, err := models.GetUserAll()
	if err != nil {
		sc.Logger.Error("获取用户列表失败", zap.Error(err))
		common.ReqBadFailWithMessage(fmt.Sprintf("获取用户列表失败:%v", err.Error()), c)
		return
	}
	resp := &ResponseResourceCommon{
		Total: len(users),
		Items: users,
	}
	common.OkWithDetailed(resp, "获取用户列表成功", c)
}

func accountExist(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	// 校验一下 menu字段
	var reqUser models.AccountExistRequest
	err := c.ShouldBindJSON(&reqUser)
	if err != nil {
		sc.Logger.Error("解析编辑用户请求失败", zap.Any("用户", reqUser), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 在这里校验字段，是否必填，范围是否正确
	err = validate.Struct(reqUser)
	if err != nil {

		// 这里为什么要判断错误是否是 ValidationErrors
		if errors, ok := err.(validator.ValidationErrors); ok {
			common.ReqBadFailWithWithDetailed(

				gin.H{
					"翻译前": err.Error(),
					"翻译后": errors.Translate(trans),
				},
				"请求出错",
				c,
			)
			return
		}
		common.ReqBadFailWithMessage(err.Error(), c)
		return

	}
	// 先 去db中根据id找到这个user

	dbUser, _ := models.GetUserByName(reqUser.Account)
	if dbUser != nil {
		sc.Logger.Info("用户已存在", zap.Any("用户", reqUser))
		common.FailWithMessage("用户已存在", c)
		return
	}
	common.OkWithMessage("用户名不存在 可用", c)
}
