package view

import (
	"cloudops/src/cache"
	"cloudops/src/common"
	"cloudops/src/config"
	"cloudops/src/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"strconv"
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

// 删除服务树节点
func deleteStreeNode(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	id := c.Param("id")
	sc.Logger.Info("删除树节点", zap.Any("id", id))

	// // 获取节点
	intVar, _ := strconv.Atoi(id)
	dbNode, err := models.GetStreeNodeById(intVar)
	if err != nil {
		sc.Logger.Error("根据id找树节点错误", zap.Any("树节点", id), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 删除权限校验
	pass, err := streeNodeOpsAdminPermissionCheck(dbNode, c)
	if err != nil {
		sc.Logger.Error("服务树节点权限校验失败", zap.Any("reqNode", dbNode), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	if !pass {
		sc.Logger.Error("服务树节点权限校验未通过", zap.Any("reqNode", dbNode), zap.Error(err))
		common.Req403WithWithMessage("服务树节点权限校验未通过", c)
		return

	}

	// // 根据dbNode的Id去查 pid
	childrens, _ := models.GetStreeNodesByPid(int(dbNode.ID))

	// 如果dbNode的children不为空 那么不允许删除
	if childrens != nil && len(childrens) > 0 {
		err = errors.New(fmt.Sprintf("不允许删除非叶子节点 id:%v title:%v",
			id,
			dbNode.Title,
		))
		sc.Logger.Error("不允许删除非叶子节点",
			zap.Error(err),
		)
		common.FailWithMessage(err.Error(), c)
		return
	}
	// 删除节点
	err = dbNode.DeleteOne()
	if err != nil {
		sc.Logger.Error("根据id删除树节点错误", zap.Any("树节点", id), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}
	common.OkWithMessage("删除成功", c)
}

// 获取所有顶级节点，为了下一轮进行子节点的展开
func getTopStreeNodes(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	sc.Logger.Info("获取顶级节点前")
	topNodes, err := models.GetStreeNodeByLevel(1)
	sc.Logger.Info("获取顶级节点后")
	if err != nil {
		sc.Logger.Error("去数据库中拿到所有的顶级服务树节点错误")
		zap.Error(err)
		common.ReqBadFailWithMessage(fmt.Sprintf("去数据库中拿到所有的顶级服务树节点错误:%v", err.Error()), c)
		return
	}

	// 遍历topNodes，获取其子节点
	for _, node := range topNodes {
		node := node
		node.FillFrontAllDataNew()
		pass, _ := streeNodeOpsAdminPermissionCheck(node, c)
		node.CanAdminNode = pass
	}
	common.OkWithDetailed(topNodes, "获取成功", c)
}

// 作用：获取所有顶级节点，为了下一轮进行子节点的展开
func getTopStreeNodesUseCache(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	sc.Logger.Info("获取顶级节点前")
	topNodes, err := models.GetStreeNodeByLevel(1)
	sc.Logger.Info("获取顶级节点后")
	if err != nil {
		sc.Logger.Error("去数据库中拿所有的顶级服务树节点错误",
			zap.Error(err),
		)
		common.ReqBadFailWithMessage(fmt.Sprintf("去数据库中拿所有的顶级服务树节点错误:%v", err.Error()), c)
		return
	}

	streeC := c.MustGet(common.GIN_CTX_STREE_CACHE).(*cache.StreeCache)
	for _, node := range topNodes {
		node := node

		node.FillFrontAllDataWithCache(streeC.StreeNodeCacahe)
		pass, _ := streeNodeOpsAdminPermissionCheck(node, c)
		node.CanAdminNode = pass
	}

	common.OkWithDetailed(topNodes, "ok", c)
}

// getChildrenStreeNodes 获取子节点
func getChildrenStreeNodes(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	pid := c.Param("pid")
	sc.Logger.Info("获取子节点", zap.Any("pid", pid))
	// 获取子节点
	streeC := c.MustGet(common.GIN_CTX_STREE_CACHE).(*cache.StreeCache)

	intVar, _ := strconv.Atoi(pid)
	childrens, err := models.GetStreeNodesByPid(intVar)
	if err != nil {
		sc.Logger.Error("根据pid获取子节点错误", zap.Any("pid", pid), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}
	for _, node := range childrens {
		node := node
		node.FillFrontAllDataWithCache(streeC.StreeNodeCacahe)
		pass, _ := streeNodeOpsAdminPermissionCheck(node, c)
		node.CanAdminNode = pass
	}
	common.OkWithDetailed(childrens, "获取成功", c)
}

// 更新服务树节点信息
// updateStreeNode 更新树节点的函数
func updateStreeNode(c *gin.Context) {
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	// 校验字段
	var reqNode models.StreeNode
	err := c.ShouldBindJSON(&reqNode)
	if err != nil {
		sc.Logger.Error("解析更新服务树节点请求失败", zap.Any("树节点", reqNode), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 校验字段，是否必填，范围是否正确
	err = validate.Struct(reqNode)
	if err != nil {

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

	_, err = models.GetStreeNodeById(int(reqNode.ID))
	if err != nil {
		sc.Logger.Error("根据id找树节点错误", zap.Any("树节点", reqNode), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	// 在创建节点前要校验权限
	pass, err := streeNodeOpsAdminPermissionCheck(&reqNode, c)
	if err != nil {
		sc.Logger.Error("服务树节点权限校验失败", zap.Any("reqNode", &reqNode), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	if !pass {
		sc.Logger.Error("服务树节点权限校验未通过", zap.Any("reqNode", &reqNode), zap.Error(err))
		common.Req403WithWithMessage("服务树节点权限校验未通过", c)
		return

	}

	usersOpsAdmin := make([]*models.User, 0)
	usersRdAdmin := make([]*models.User, 0)
	usersRdMember := make([]*models.User, 0)
	// 遍历角色menu 列表 找到角色
	for _, userName := range reqNode.OpsAdminUsers {
		dbUser, err := models.GetUserByUserName(userName)
		if err != nil {
			sc.Logger.Error("树节点根据userName找用户错误", zap.Any("树节点", reqNode), zap.Error(err))
			common.FailWithMessage(err.Error(), c)
			return
		}
		usersOpsAdmin = append(usersOpsAdmin, dbUser)

	}
	for _, userName := range reqNode.RdAdminUsers {
		dbUser, err := models.GetUserByUserName(userName)
		if err != nil {
			sc.Logger.Error("树节点根据userName找用户错误", zap.Any("树节点", reqNode), zap.Error(err))
			common.FailWithMessage(err.Error(), c)
			return
		}
		usersRdAdmin = append(usersRdAdmin, dbUser)

	}
	for _, userName := range reqNode.RdMemberUsers {
		dbUser, err := models.GetUserByUserName(userName)
		if err != nil {
			sc.Logger.Error("树节点根据userName找用户错误", zap.Any("树节点", reqNode), zap.Error(err))
			common.FailWithMessage(err.Error(), c)
			return
		}
		usersRdMember = append(usersRdMember, dbUser)

	}

	// 将运维负责人设置，即使是空的也要设置：
	// 原来非空，现在为空 代表删除人员了
	reqNode.OpsAdmins = usersOpsAdmin
	reqNode.RdAdmins = usersRdAdmin
	reqNode.RdMembers = usersRdMember
	err = reqNode.UpdateStreeNode()
	if err != nil {
		sc.Logger.Error("更新树节点和关联的运维负责人错误", zap.Any("树节点", reqNode), zap.Error(err))
		common.FailWithMessage(err.Error(), c)
		return
	}

	common.OkWithMessage("更新成功", c)
	sc.Logger.Info("更新树节点和关联的运维负责人成功", zap.Any("树节点", reqNode))
}
