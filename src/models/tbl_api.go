package models

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Api 结构体定义了API的数据库表字段及其关系
type Api struct {
	Model            // 嵌入 Model 结构体，自动包含 ID 和 CreatedAt 字段
	Path     string  `json:"path" gorm:"type:varchar(50);comment:路由路径"`            // API 路径
	Method   string  `json:"method" gorm:"type:varchar(50);comment:http请求方法"`      // HTTP 请求方法
	Pid      int     `json:"pId" gorm:"comment:apiGroups 父级的id 为了给树用的"`            // 父级 API ID，用于树形结构
	Title    string  `json:"title" gorm:"type:varchar(50);uniqueIndex;comment:名称"` // API 名称，唯一索引
	Roles    []*Role `json:"roles" gorm:"many2many:role_apis;"`                    // 多对多关系，API 对应的角色
	Type     string  `json:"type" gorm:"type:varchar(5);comment:类型 0=父级 1=子级"`     // API 类型
	Key      uint    `json:"key"  gorm:"-"`                                        // 前端用的键
	Value    uint    `json:"value"  gorm:"-"`                                      // 前端用的值
	Children []*Api  `json:"children" gorm:"-"`                                    // 子 API，返回给前端
}

// Create 创建一个新api记录
func (obj *Api) Create() error {
	return DB.Create(obj).Error
}

// GetApiAll 获取所有API记录
func GetApiAll() (objs []*Api, err error) {
	err = DB.Find(&objs).Error
	return
}

// GetApiById 根据 ID 获取 API 包括关联的角色信息
func GetApiById(id int) (*Api, error) {
	var dbObje Api
	// 根据ID查询API，并预加载关联的角色信息
	err := DB.Where("id = ?").Preload("Roles").First(&dbObje).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("api不存在")
		}
		return nil, fmt.Errorf("数据库错误:%w", err)
	}
	return &dbObje, nil
}

// GetApiByTitle 根据标题获取API记录，包括关联的角色信息
func GetApiByTitle(title string) (*Api, error) {
	var dbObje Api
	// 根据标题查询API，并预加载关联的角色信息
	err := DB.Where("title = ?", title).Preload("Roles").First(&dbObje).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("api不存在")
		}
		return nil, fmt.Errorf("数据库错误:%w", err)
	}
	return &dbObje, nil
}

// DeleteOne 删除一条API记录，包扩其所有关联
func (obj *Api) DeleteOne() error {
	return DB.Select(clause.Associations).Unscoped().Delete(obj).Error
}

// CreateOne 创建一个新API记录
func (obj *Api) CreateOne() error {
	return DB.Create(obj).Error
}

// UpdateOne 更新一条API记录
func (obj *Api) UpdateOne() error {
	return DB.Save(obj).Error
}
