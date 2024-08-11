package models

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StreeNode struct {
	Model             // 不用每次写ID 和 createAt了
	Title     string  `json:"title" gorm:"uniqueIndex:pid_title;type:varchar(50);comment:名称" `
	Pid       uint    `json:"pId" gorm:"index;uniqueIndex:pid_title;comment:StreeNodeGroups 父级的id 为了给树用的"`
	Level     int     `json:"level" gorm:"comment:层级"`
	IsLeaf    bool    `json:"isLeaf" gorm:"comment:是否启用 0=否 1=是"`
	Desc      string  `json:"desc"  gorm:"comment:描述"`
	OpsAdmins []*User `json:"ops_admins" gorm:"many2many:ops_admins;comment:运维负责人列表 可以做运维操作"`
	RdAdmins  []*User `json:"rd_admins" gorm:"many2many:rd_admins;comment:研发负责人列表 可以审批发布单"`          // 多对多
	RdMembers []*User `json:"rd_members" gorm:"many2many:rd_members;comment:研发工程师列表 可以提发布单 可以操作发布单"` // 多对多

	Children []*StreeNode `json:"children" gorm:"-"` // 返回给前端使用

	Key string `json:"key" gorm:"-"` // 返回给前端表格
}

// 创建服务树节点
func (obj *StreeNode) Creat() error {
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
	//err := DB.Where("id = ?", id).Preload("BindEcss").Preload("BindElbs").Preload("BindRdss").Preload("OpsAdmins").Preload("RdAdmins").Preload("RdMembers").First(&dbObj).Error
	err := DB.Where("id = ?", id).First(&dbObj).Error
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
	//err = DB.Where("pid = ?", pid).Preload("BindEcss").Preload("BindElbs").Preload("BindRdss").Preload("OpsAdmins").Preload("RdAdmins").Preload("RdMembers").Find(&dbObjs).Error
	err = DB.Where("pid = ?", pid).Find(&dbObjs).Error
	return
}
