package models

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"sort"
	"sync"
)

// 服务树配置结构体
type StreeNode struct {
	Model         // 不用每次写ID 和 createAt了
	Title  string `json:"title" gorm:"uniqueIndex:pid_title;type:varchar(50);comment:名称" `
	Pid    uint   `json:"pId" gorm:"index;uniqueIndex:pid_title;comment:StreeNodeGroups 父级的id 为了给树用的"`
	Level  int    `json:"level" gorm:"comment:层级 "`
	IsLeaf bool   `json:"isLeaf" gorm:"comment:是否启用 0=否 1=是"`
	Desc   string `json:"desc"  gorm:"comment:描述"`

	OpsAdmins []*User `json:"ops_admins" gorm:"many2many:ops_admins;comment:运维负责人列表 可以做运维操作"`        // 多对多
	RdAdmins  []*User `json:"rd_admins" gorm:"many2many:rd_admins;comment:研发负责人列表 可以审批发布单"`          // 多对多
	RdMembers []*User `json:"rd_members" gorm:"many2many:rd_members;comment:研发工程师列表 可以提发布单 可以操作发布单"` // 多对多

	BindEcss         []*ResourceEcs `json:"bind_ecss" gorm:"many2many:bind_ecss;"` // 多对多
	BindElbs         []*ResourceElb `json:"bind_elbs" gorm:"many2many:bind_elbs;"` // 多对多
	BindRdss         []*ResourceRds `json:"bind_rdss" gorm:"many2many:bind_rdss;"` // 多对多
	OpsAdminUsers    []string       `json:"ops_admin_users" gorm:"-"`              // 节点运维负责人名字列表
	RdAdminUsers     []string       `json:"rd_admin_users" gorm:"-"`               // 节点研发负责人名字列表
	RdMemberUsers    []string       `json:"rd_member_users" gorm:"-"`              // 节点研发工程师名字列表
	EcsNum           int            `json:"ecsNum" gorm:"-"`
	ElbNum           int            `json:"elbNum" gorm:"-"`
	RdsNum           int            `json:"rdsNum" gorm:"-"`
	NodeNum          int            `json:"nodeNum" gorm:"-"`          // 子节点数量
	LeafNodeNum      int            `json:"leafNodeNum" gorm:"-"`      // 叶子节点数量
	EcsCpuTotal      int            `json:"ecsCpuTotal" gorm:"-"`      // 叶子节点数量
	ElbBandWithTotal int            `json:"elbBandWithTotal" gorm:"-"` // elb带宽包资源上限
	EcsMemoryTotal   int            `json:"ecsMemoryTotal" gorm:"-"`   // 叶子节点数量
	EcsDiskTotal     int            `json:"ecsDiskTotal" gorm:"-"`     // 叶子节点数量
	Children         []*StreeNode   `json:"children" gorm:"-"`         // 返回给前端的
	//Key              uint           `json:"key" gorm:"-"`              // 返回给前端的
	Key   string `json:"key" gorm:"-"`   //给前端表格
	Label string `json:"label" gorm:"-"` //给前端表格
	Value uint   `json:"value" gorm:"-"` // 返回给前端的

	OpsRdAdmins []string `json:"opsRdAdmins" gorm:"-"`

	GroupByVendor            []*EchartOneItem `json:"groupByVendor" gorm:"-"`
	GroupByVendorMap         map[string]int   `json:"-" gorm:"-"`
	GroupByZoneIdMap         map[string]int   `json:"-" gorm:"-"`
	GroupByOSName            map[string]int   `json:"-" gorm:"-"`
	GroupByZoneId            []*EchartOneItem `json:"groupByZoneId" gorm:"-"`
	GroupByOSNameOrderKeys   []string         `json:"groupByOSNameOrderKeys" gorm:"-"`
	GroupByOSNameOrderValues []int            `json:"groupByOSNameOrderValues" gorm:"-"`

	// elb
	GroupByVendorElb              []*EchartOneItem `json:"groupByVendorElb" gorm:"-"`
	GroupByZoneIdElb              []*EchartOneItem `json:"groupByZoneIdElb" gorm:"-"`
	GroupByLoadBalancerTypeKeys   []string         `json:"groupByLoadBalancerTypeKeys" gorm:"-"`
	GroupByLoadBalancerTypeValues []int            `json:"groupByLoadBalancerTypeValues" gorm:"-"`

	// rds
	GroupByVendorRds       []*EchartOneItem `json:"groupByVendorRds" gorm:"-"`
	GroupByZoneIdRds       []*EchartOneItem `json:"groupByZoneIdRds" gorm:"-"`
	GroupByRdsEngineKeys   []string         `json:"groupByRdsEngineKeys" gorm:"-"`
	GroupByRdsEngineValues []int            `json:"groupByRdsEngineValues" gorm:"-"`

	// 是否有权限操作这节点
	CanAdminNode bool   `json:"canAdminNode" gorm:"-"`
	NodePath     string `json:"nodePath" gorm:"-"` // a.b.c.d
}

// 创建服务树节点
func (obj *StreeNode) Create() error {
	return DB.Create(obj).Error

}

// 删除服务树节点
func (obj *StreeNode) DeleteOne() error {
	return DB.Select(clause.Associations).Unscoped().Delete(obj).Error
}

// 创建一个新的StreeNode对象
func (obj *StreeNode) CreateOne() error {
	return DB.Create(obj).Error
}

// 更新StreeNode对象
func (obj *StreeNode) UpdateOne() error {
	return DB.Updates(obj).Error

}

// 获取所有StreeNode对象
func GetStreeNodeAll() (objs []*StreeNode, err error) {
	err = DB.Find(&objs).Error
	return
}

// 查找服务树节点ID
func GetStreeNodeById(id int) (*StreeNode, error) {
	var dbObj StreeNode
	err := DB.Where("id = ?", id).Preload("BindEcss").Preload("BindElbs").Preload("BindRdss").Preload("OpsAdmins").Preload("RdAdmins").Preload("RdMembers").First(&dbObj).Error
	//err := DB.Where("id = ?", id).Preload("OpsAdmins").Preload("RdAdmins").Preload("RdMembers").First(&dbObj).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("streeNode不存在")
		}
		return nil, fmt.Errorf("数据库错误:%w", err)
	}
	return &dbObj, nil

}

// 查找服务树节点Pid
func GetStreeNodesByPid(pid int) (dbObjs []*StreeNode, err error) {
	err = DB.Where("pid = ?", pid).Preload("BindEcss").Preload("BindElbs").Preload("BindRdss").Preload("OpsAdmins").Preload("RdAdmins").Preload("RdMembers").Find(&dbObjs).Error
	//err = DB.Where("pid = ?", pid).Find(&dbObjs).Error
	return
}

// 获取服务树节点层级
func GetStreeNodeByLevel(level int) (objs []*StreeNode, err error) {
	err = DB.Debug().Where("level = ?", level).Preload("BindEcss").Preload("BindElbs").Preload("BindRdss").Preload("OpsAdmins").Preload("RdAdmins").Preload("RdMembers").Find(&objs).Error
	return
}

func GetStreeNodesByPidNew(pid int) (dbObjs []*StreeNode, err error) {
	err = DB.Debug().Where("pid = ?", pid).Preload("BindEcss").Preload("BindElbs").Preload("BindRdss").Preload("OpsAdmins").Preload("RdAdmins").Preload("RdMembers").Find(&dbObjs).Error
	return
}

// 不带preload
func GetStreeNodesByPidNoPreload(pid int) (dbObjs []*StreeNode, err error) {
	err = DB.Debug().Where("pid = ?", pid).Find(&dbObjs).Error
	return
}

func (obj *StreeNode) FillFrontAllDataNew() {
	//obj.Key = obj.ID
	//obj.Key = obj.ID

	obj.Key = fmt.Sprintf("%d", obj.ID)
	log.Printf("FillFrontAllData start FillUsers node:%v", obj.Title)
	obj.FillUsers()
	log.Printf("FillFrontAllData end FillUsers node:%v", obj.Title)
	log.Printf("FillFrontAllData start GetFullNodePath node:%v", obj.Title)
	obj.GetFullNodePath()
	log.Printf("FillFrontAllData end GetFullNodePath node:%v", obj.Title)
	log.Printf("FillFrontAllData start FillFrontResource node:%v", obj.Title)
	obj.FillFrontResource()
	log.Printf("FillFrontAllData end FillFrontResource node:%v", obj.Title)
	log.Printf("FillFrontAllData start BindEcsData node:%v", obj.Title)
	obj.BindEcsData()
	log.Printf("FillFrontAllData end BindEcsData node:%v", obj.Title)
	log.Printf("FillFrontAllData start BindElbData node:%v", obj.Title)
	obj.BindElbData()
	log.Printf("FillFrontAllData end BindElbData node:%v", obj.Title)
	log.Printf("FillFrontAllData start BindRdsData node:%v", obj.Title)
	obj.BindRdsData()
	log.Printf("FillFrontAllData end BindRdsData node:%v", obj.Title)

	//obj.SetEcsNum()
}

// 填充人员信息
func (obj *StreeNode) FillUsers() {
	obj.OpsAdminUsers = []string{}
	obj.RdAdminUsers = []string{}
	obj.RdMemberUsers = []string{}
	for _, user := range obj.OpsAdmins {
		user := user
		obj.OpsAdminUsers = append(obj.OpsAdminUsers, user.Username)
	}
	for _, user := range obj.RdAdmins {
		user := user
		obj.RdAdminUsers = append(obj.RdAdminUsers, user.Username)
	}
	for _, user := range obj.RdMembers {
		user := user
		obj.RdMemberUsers = append(obj.RdMemberUsers, user.Username)
	}
	cicdCanApprovalAll := obj.GetAllRdOpsAdmins()
	obj.OpsRdAdmins = []string{}
	for _, user := range cicdCanApprovalAll {
		user := user
		obj.OpsRdAdmins = append(obj.OpsRdAdmins, user.Username)
	}

}

// 获取全路径
func (obj *StreeNode) GetFullNodePath() error {

	// 不断的遍历这个node 的pid
	fatherTitles := []string{}

	node := obj
	for node.Pid > 0 {
		father, err := GetStreeNodeById(int(node.Pid))
		if err != nil {
			return err
		}
		fatherTitles = append(fatherTitles, father.Title)
		node = father
	}
	// 这时 fatherTitles 就是倒着的
	nodePath := ""
	num := len(fatherTitles)
	for i := num - 1; i >= 0; i-- {
		title := fatherTitles[i]
		if nodePath == "" {
			nodePath = title
			continue
		}
		nodePath = fmt.Sprintf("%s.%s", nodePath, title)
	}
	if nodePath != "" {
		nodePath = fmt.Sprintf("%s.%s", nodePath, obj.Title)
	} else {
		nodePath = obj.Title
	}

	obj.NodePath = nodePath
	return nil

}

// 获取CICD工单审批人
func (obj *StreeNode) GetAllRdOpsAdmins() []*User {
	all := []*User{}
	// 本层级
	if obj.RdAdmins != nil && len(obj.RdAdmins) > 0 {
		all = append(all, obj.RdAdmins...)
	}
	if obj.OpsAdmins != nil && len(obj.OpsAdmins) > 0 {
		all = append(all, obj.OpsAdmins...)
	}
	if len(all) > 0 {
		return all
	}
	if obj.Pid > 0 {
		father, _ := GetStreeNodeById(int(obj.Pid))
		return father.GetAllRdOpsAdmins()
	} else {
		return nil
	}
}

// 填充前端所需的字段
// 绑定资源key为了穿梭框
func (obj *StreeNode) FillFrontResource() {

	for _, obj := range obj.BindEcss {
		obj := obj
		obj.Key = fmt.Sprintf("%d", obj.ID)
	}

	for _, obj := range obj.BindElbs {
		obj := obj
		obj.Key = fmt.Sprintf("%d", obj.ID)
	}

	for _, obj := range obj.BindRdss {
		obj := obj
		obj.Key = fmt.Sprintf("%d", obj.ID)
	}
}

func GetAllLeafNodesNew(pid int) (objs []*StreeNode) {
	childrens, _ := GetStreeNodesByPidNew(pid)
	if len(childrens) == 0 {
		return
	}
	// 添加结果
	objs = append(objs, childrens...)
	// 遍历childrens 做递归
	for _, node := range childrens {
		node := node
		objs = append(objs, GetAllLeafNodesNew(int(node.ID))...)
	}
	return
}

// Ecs统计计算
func (obj *StreeNode) BindEcsData() {
	allNum := 0
	groupByVendor := make(map[string]int)
	groupByOSName := make(map[string]int)
	groupByZoneId := make(map[string]int)
	allNodes := []*StreeNode{obj}

	// 对比拼接的速度
	allNodes = append(allNodes, GetAllLeafNodesNew(int(obj.ID))...)
	//allNodes = append(allNodes, GetAllLeafNodesPinJie(int(obj.ID))...)

	// 去掉自身算子节点数量
	obj.NodeNum = len(allNodes) - 1
	allResourcesIdsMap := map[uint]struct{}{}

	for _, node := range allNodes {
		if node.IsLeaf {
			obj.LeafNodeNum++
		}
		//fmt.Printf("[father:%v allNodes-num:%v]BindEcsData .node.print : index:% node %v\n", obj.Title, len(allNodes), index, node)
		node := node
		if node.BindEcss == nil {
			continue
		}
		for _, objEcs := range node.BindEcss {
			objEcs := objEcs
			groupByVendor[objEcs.Vendor]++
			groupByOSName[objEcs.OSName]++
			groupByZoneId[objEcs.ZoneId]++
			obj.EcsCpuTotal += objEcs.Cpu
			obj.EcsMemoryTotal += objEcs.Memory
			obj.EcsDiskTotal += objEcs.Disk

			allResourcesIdsMap[objEcs.ID] = struct{}{}
		}

	}

	allNum = len(allResourcesIdsMap)
	// 赋值

	obj.EcsNum = allNum
	arrGroupByVendor := make([]*EchartOneItem, 0)
	arrGroupByZoneId := make([]*EchartOneItem, 0)
	for name, value := range groupByVendor {
		arrGroupByVendor = append(arrGroupByVendor,
			&EchartOneItem{
				Name:  name,
				Value: value,
			},
		)
	}

	for name, value := range groupByZoneId {
		arrGroupByZoneId = append(arrGroupByZoneId,
			&EchartOneItem{
				Name:  name,
				Value: value,
			},
		)
	}

	// 对于osName 排序。先排序key

	osNameKeys := []string{}
	osNameValues := []int{}
	for name := range groupByOSName {

		osNameKeys = append(osNameKeys, name)

	}

	sort.Strings(osNameKeys)
	for _, name := range osNameKeys {
		osNameValues = append(osNameValues, groupByOSName[name])
	}

	obj.GroupByVendor = arrGroupByVendor
	obj.GroupByZoneId = arrGroupByZoneId
	obj.GroupByOSNameOrderKeys = osNameKeys
	obj.GroupByOSNameOrderValues = osNameValues
	return
}

// 使用缓存递归获取所有孩子
func (obj *StreeNode) StatisticsRecursionWithCache(p sync.Map, allNodes []*StreeNode) {
	// 遍历手动求和
	GroupByVendorMap := map[string]int{}
	GroupByZoneIdMap := map[string]int{}
	GroupByOSName := map[string]int{}
	for _, node := range allNodes {
		node := node
		key := fmt.Sprintf("Rdss-%d", node.ID)
		cacheNodeobj, ok := p.Load(key)
		if !ok {
			continue
		}
		caheNodeobj := cacheNodeobj.(*StreeNode)
		obj.EcsNum += caheNodeobj.EcsNum
		obj.ElbNum += caheNodeobj.ElbNum
		obj.RdsNum += caheNodeobj.RdsNum
		obj.LeafNodeNum += caheNodeobj.LeafNodeNum
		obj.EcsCpuTotal += caheNodeobj.EcsCpuTotal
		obj.ElbBandWithTotal += caheNodeobj.ElbBandWithTotal
		obj.EcsMemoryTotal += caheNodeobj.EcsMemoryTotal
		obj.EcsDiskTotal += caheNodeobj.EcsDiskTotal
		for k, v := range caheNodeobj.GroupByVendorMap {
			GroupByVendorMap[k] += v
		}
		for k, v := range caheNodeobj.GroupByZoneIdMap {
			GroupByZoneIdMap[k] += v
		}
		for k, v := range caheNodeobj.GroupByOSName {
			GroupByOSName[k] += v
		}
	}
	obj.GroupByVendorMap = GroupByVendorMap
	obj.GroupByZoneIdMap = GroupByZoneIdMap
	obj.GroupByOSName = GroupByOSName
	obj.Statistics()
}

// 统计计算
func (obj *StreeNode) Statistics() {
	arrGroupByVendor := make([]*EchartOneItem, 0)
	arrGroupByZoneId := make([]*EchartOneItem, 0)
	for name, value := range obj.GroupByVendorMap {
		arrGroupByVendor = append(arrGroupByVendor,
			&EchartOneItem{
				Name:  name,
				Value: value,
			},
		)
	}
	for name, value := range obj.GroupByZoneIdMap {
		arrGroupByZoneId = append(arrGroupByZoneId,
			&EchartOneItem{
				Name:  name,
				Value: value,
			},
		)
	}
	osNameKeys := []string{}
	osNameValues := []int{}
	for name := range obj.GroupByOSName {
		osNameKeys = append(osNameKeys, name)
	}
	sort.Strings(osNameKeys)
	for _, name := range osNameKeys {
		osNameValues = append(osNameValues, obj.GroupByOSName[name])
	}
	obj.GroupByVendor = arrGroupByVendor
	obj.GroupByZoneId = arrGroupByZoneId
	obj.GroupByOSNameOrderKeys = osNameKeys
	obj.GroupByOSNameOrderValues = osNameValues
}

// 目的是获取所有子孙中的叶子节点：因为只有叶子节点才可以绑定资源
func GetAllLeafNodes(pid int) (objs []*StreeNode) {
	childrens, _ := GetStreeNodesByPid(pid)
	if len(childrens) == 0 {
		return
	}
	// 添加结果
	objs = append(objs, childrens...)
	// 遍历childrens 做递归
	for _, node := range childrens {
		node := node
		objs = append(objs, GetAllLeafNodes(int(node.ID))...)
	}
	return
}

// ELB统计计算
func (obj *StreeNode) BindElbData() {
	allNum := 0
	groupByVendor := make(map[string]int)
	groupByLoadBalancerType := make(map[string]int)
	groupByZoneId := make(map[string]int)
	allNodes := []*StreeNode{obj}

	allNodes = append(allNodes, GetAllLeafNodes(int(obj.ID))...)

	// 去掉自身算子节点数量
	obj.NodeNum = len(allNodes) - 1
	allResourcesIdsMap := map[uint]struct{}{}

	for _, node := range allNodes {
		//fmt.Printf("[father:%v allNodes-num:%v]BindEcsData .node.print : index:% node %v\n", obj.Title, len(allNodes), index, node)
		node := node
		if node.BindElbs == nil {
			continue
		}
		for _, objOne := range node.BindElbs {
			objOne := objOne
			groupByVendor[objOne.Vendor]++
			groupByLoadBalancerType[objOne.LoadBalancerType]++
			groupByZoneId[objOne.ZoneId]++
			obj.ElbBandWithTotal += objOne.BandwidthCapacity

			allResourcesIdsMap[objOne.ID] = struct{}{}
		}

	}

	allNum = len(allResourcesIdsMap)
	// 赋值

	obj.ElbNum = allNum
	arrGroupByVendor := make([]*EchartOneItem, 0)
	arrGroupByZoneId := make([]*EchartOneItem, 0)
	for name, value := range groupByVendor {
		arrGroupByVendor = append(arrGroupByVendor,
			&EchartOneItem{
				Name:  name,
				Value: value,
			},
		)
	}

	for name, value := range groupByZoneId {
		arrGroupByZoneId = append(arrGroupByZoneId,
			&EchartOneItem{
				Name:  name,
				Value: value,
			},
		)
	}

	// 对于osName 排序。先排序key

	elbTypeNameKeys := []string{}
	elbTypeValues := []int{}
	for name := range groupByLoadBalancerType {

		elbTypeNameKeys = append(elbTypeNameKeys, name)

	}

	sort.Strings(elbTypeNameKeys)
	for _, name := range elbTypeNameKeys {
		elbTypeValues = append(elbTypeValues, groupByLoadBalancerType[name])
	}

	obj.GroupByVendorElb = arrGroupByVendor
	obj.GroupByZoneIdElb = arrGroupByZoneId
	obj.GroupByLoadBalancerTypeKeys = elbTypeNameKeys
	obj.GroupByLoadBalancerTypeValues = elbTypeValues
	return
}

// Rds统计计算
func (obj *StreeNode) BindRdsData() {
	allNum := 0
	groupByVendor := make(map[string]int)
	groupByEngine := make(map[string]int)
	groupByZoneId := make(map[string]int)
	allNodes := []*StreeNode{obj}

	allNodes = append(allNodes, GetAllLeafNodes(int(obj.ID))...)

	// 去掉自身算子节点数量
	obj.NodeNum = len(allNodes) - 1
	allResourcesIdsMap := map[uint]struct{}{}

	for _, node := range allNodes {

		//fmt.Printf("[father:%v allNodes-num:%v]BindEcsData .node.print : index:% node %v\n", obj.Title, len(allNodes), index, node)
		node := node
		if node.BindRdss == nil {
			continue
		}
		for _, objOne := range node.BindRdss {
			objOne := objOne
			groupByVendor[objOne.Vendor]++
			groupByEngine[objOne.Engine]++
			groupByZoneId[objOne.ZoneId]++

			allResourcesIdsMap[objOne.ID] = struct{}{}
		}

	}

	allNum = len(allResourcesIdsMap)
	// 赋值

	obj.RdsNum = allNum
	arrGroupByVendor := make([]*EchartOneItem, 0)
	arrGroupByZoneId := make([]*EchartOneItem, 0)
	for name, value := range groupByVendor {
		arrGroupByVendor = append(arrGroupByVendor,
			&EchartOneItem{
				Name:  name,
				Value: value,
			},
		)
	}

	for name, value := range groupByZoneId {
		arrGroupByZoneId = append(arrGroupByZoneId,
			&EchartOneItem{
				Name:  name,
				Value: value,
			},
		)
	}

	// 对于osName 排序。先排序key

	rdsEngineKeys := []string{}
	rdsEngineValues := []int{}
	for name := range groupByEngine {

		rdsEngineKeys = append(rdsEngineKeys, name)

	}

	sort.Strings(rdsEngineKeys)
	for _, name := range rdsEngineKeys {
		rdsEngineValues = append(rdsEngineValues, groupByEngine[name])
	}

	obj.GroupByVendorRds = arrGroupByVendor
	obj.GroupByZoneIdRds = arrGroupByZoneId
	obj.GroupByRdsEngineKeys = rdsEngineKeys
	obj.GroupByRdsEngineValues = rdsEngineValues
	return
}

// GetAllLeafNodesNoPreload
func GetAllLeafNodesNoPreload(pid int) (objs []*StreeNode) {
	childrens, _ := GetStreeNodesByPidNoPreload(pid)
	if len(childrens) == 0 {
		return
	}
	// 添加结果
	objs = append(objs, childrens...)
	// 遍历chuldrens
	for _, node := range childrens {
		objs = append(objs, GetAllLeafNodesNoPreload(int(node.ID))...)
	}
	return
}

// FillFrontAllDataWithCache
func (obj *StreeNode) FillFrontAllDataWithCache(p sync.Map) {

	obj.Key = fmt.Sprintf("%d", obj.ID)
	obj.FillUsers()
	obj.GetFullNodePath()
	obj.FillFrontResource()

	// 递归获取所有子节点
	allNodes := []*StreeNode{obj}
	allNodes = append(allNodes, GetAllLeafNodesNoPreload(int(obj.ID))...)

	obj.StatisticsRecursionWithCache(p, allNodes)
}

// 更新服务树节点信息

func (obj *StreeNode) UpdateStreeNode() error {
	//log.Printf("准备更新树节点: %+v", obj)
	// 更新节点本身
	err := DB.Where("id = ?", obj.ID).Updates(obj).Error
	if err != nil {
		return fmt.Errorf("更新本体错误: %w", err)
	}

	// 更新关联的运维负责人
	err = DB.Model(obj).Association("OpsAdmins").Replace(obj.OpsAdmins)
	if err != nil {
		return fmt.Errorf("更新关联运维负责人错误: %w", err)
	}

	// 更新关联的研发负责人
	err = DB.Model(obj).Association("RdAdmins").Replace(obj.RdAdmins)
	if err != nil {
		return fmt.Errorf("更新关联研发负责人错误: %w", err)
	}

	// 更新关联的研发工程师
	err = DB.Model(obj).Association("RdMembers").Replace(obj.RdMembers)
	if err != nil {
		return fmt.Errorf("更新关联研发工程师错误: %w", err)
	}

	return nil // 如果没有错误发生，返回 nil
}
