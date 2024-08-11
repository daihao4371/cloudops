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

type tmpNode struct {
	Title    string     `json:"title"`
	Key      string     `json:"key"`
	Children []*tmpNode `json:"children,omitempty"`
	Pid      int        `json:"pid"`
	Level    int        `json:"level"`
}

// 获取服务树所有列表
func getStreeNodeList(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	//streeNodes, err := models.GetStreeNodeAll()
	streeNodes, err := models.GetStreeNodeAll()
	if err != nil {
		sc.Logger.Error("去数据库中拿到所有的服务树节点错误", zap.Error(err))
		common.ReqBadFailWithMessage(fmt.Sprintf("去数据库拿到所有的服务树节点错误:%s", err.Error()), c)
		return
	}

	// 准备全员ID map
	allMap := map[uint]*models.StreeNode{}
	topMap := map[uint]*models.StreeNode{}
	pidMap := map[uint][]*models.StreeNode{}

	// 遍历所有节点，将所有节点放入map中
	for _, streeNode := range streeNodes {
		streeNode := streeNode
		streeNode.Key = fmt.Sprintf("%d", streeNode.ID)

		allMap[streeNode.ID] = streeNode
		if streeNode.Pid == 0 {
			topMap[streeNode.ID] = streeNode
			continue
		}
		//同一层级放到一个map中
		chidls, ok := pidMap[streeNode.Pid]
		if !ok {
			chidls = []*models.StreeNode{}
		}
		chidls = append(chidls, streeNode)
		pidMap[streeNode.Pid] = chidls
	}

	// 再遍历allMap 回填他们的 children 列表
	for id, node := range allMap {
		id := id
		node := node
		childs, ok := pidMap[id]
		if ok {
			node.Children = childs
		}
	}
	finalNodes := []*models.StreeNode{}

	// 遍历topMap，将顶级节点放入finalNodes中
	for _, node := range topMap {
		node := node
		finalNodes = append(finalNodes, node)
	}
	common.OkWithDetailed(finalNodes, "ok", c)
}

// 增删改查树节点的通用权限校验方法 true 代表有权限 fale 代表没有权限
func streeNodeOpsAdminPermissionCheck(node *models.StreeNode, c *gin.Context) (bool, error) {
	userName := c.MustGet(common.GIN_CTX_JWT_USER_NAME).(string)
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	reqUser, err := models.GetUserByUserName(userName)

	if err != nil {
		sc.Logger.Error("通过token解析到的userName去数据库中查找用户错误", zap.Error(err))
		common.ReqBadFailWithMessage(fmt.Sprintf("通过token解析到的userName去数据库中查找用户错误:%s", err.Error()), c)
		return false, err
	}

	// 判断用户是否是管理员
	if reqUser.AccountType == 2 {
		return true, nil
	}

	sc.Logger.Debug("服务树节点权限校验开始",
		zap.String("用户", userName),
		zap.Any("node", node),
	)

	// 如果是超管就放行
	isSuper := false
	for _, role := range reqUser.Roles {
		role := role
		if role.RoleValue == "super" {
			isSuper = true
			break
		}
	}

	if isSuper {
		return true, nil
	}
	pass := false

	// 否则的话，校验用户在不在，运维管理员列表里面，或者父级
	for node != nil {
		// 先遍历OpsAdmins
		for _, user := range node.OpsAdmins {
			user := user
			if user.Username == reqUser.Username {
				pass = true
				break
			}
		}
		if pass {
			break
		}
		// 父级
		father, _ := models.GetStreeNodeById(int(node.Pid))
		if father == nil {
			break
		}
		node = father
	}
	sc.Logger.Info("服务树节点权限校验结束",
		zap.Bool("是否通过", pass),
		zap.String("用户", userName),
	)
	return pass, nil
}

// 创建服务树节点
func createStreeNode(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	// 字段校验
	var reqNode models.StreeNode
	err := c.ShouldBindJSON(&reqNode)
	if err != nil {
		sc.Logger.Error("解析新增reqNode请求失败，", zap.Any("reqNode", reqNode), zap.Error(err))
		common.FailWithMessage(fmt.Sprintf("请求参数错误:%s", err.Error()), c)
		return
	}

	// 字段校验，是否必填，范围是否正确
	err = validate.Struct(reqNode)
	if err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			common.ReqBadFailWithWithDetailed(gin.H{
				"翻译前": err.Error(),
				"翻译后": errors.Translate(trans),
			},
				"请求出错",
				c,
			)
			return
		}
		common.ReqBadFailWithMessage(err.Error(), c)
	}

	// 创建节点前要校验权限
	pass, err := streeNodeOpsAdminPermissionCheck(&reqNode, c)
	if err != nil {
		sc.Logger.Error("服务树节点权限校验失败", zap.Any("reqNode", reqNode), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}
	if !pass {
		sc.Logger.Error("服务树节点权限校验失败", zap.Any("reqNode", reqNode), zap.Error(err))
		common.FailWithMessage("服务树节点权限校验未通过", c)
		return
	}
	err = reqNode.CreateOne()
	if err != nil {
		sc.Logger.Error("创建服务树节点失败", zap.Any("reqNode", reqNode), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}
	common.OkWithMessage("创建成功", c)
}
