package models

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Menu 结构体定义了菜单的数据库表字段及其关系
type Menu struct {
	gorm.Model           // 嵌入 Model 结构体，自动包含 ID 和 CreatedAt 字段
	Name       string    `json:"name" gorm:"type:varchar(100);uniqueIndex;comment:名称"`      // 菜单名称，唯一索引
	Title      string    `json:"title" gorm:"comment:名称"`                                   // 菜单标题
	Pid        string    `json:"pid" gorm:"comment:父级的id"`                                  // 父级菜单ID
	ParentMenu string    `json:"parentMenu" gorm:"varchar(5);comment:父级的id"`                // 父级菜单ID，字符串形式
	Meta       *MenuMeta `json:"meta" gorm:"-"`                                             // 菜单元数据，不保存到数据库
	Icon       string    `json:"icon" gorm:"comment:图标"`                                    // 图标
	Type       string    `json:"type" gorm:"type:varchar(5);comment:类型 0=目录 1=子菜单"`         // 菜单类型
	Show       string    `json:"show" gorm:"type:varchar(5);comment:类型 0=禁用 1=启用"`          // 是否显示菜单
	OrderNo    int       `json:"orderNo" gorm:"comment:排序"`                                 // 排序号
	Component  string    `json:"component" gorm:"type:varchar(50);comment:前端组件 菜单就是LAYOUT"` // 组件名称
	Redirect   string    `json:"redirect" gorm:"type:varchar(50);comment:显示路径"`             // 重定向路径
	Path       string    `json:"path" gorm:"type:varchar(50);comment:路由路径"`                 // 路由路径
	Remark     string    `json:"remark" gorm:"comment:用户描述"`                                // 描述
	HomePath   string    `json:"homePath" gorm:"comment:登陆后的默认首页"`                          // 登录后的默认首页
	Status     string    `json:"status" gorm:"default:1;comment:是否启用 0禁用 1启用"`              // 菜单状态
	Roles      []*Role   `json:"roles" gorm:"many2many:role_menus;"`                        // 多对多关系，菜单对应的角色
	Children   []*Menu   `json:"children" gorm:"-"`                                         // 子菜单，返回给前端
	Key        uint      `json:"value"  gorm:"-"`                                           // 前端用的键
	Value      uint      `json:"key"  gorm:"-"`                                             // 前端用的值
}

// MenuMeta 结构体定义了菜单的元数据字段
type MenuMeta struct {
	Title           string `json:"title" gorm:"-"`           // 菜单标题
	Icon            string `json:"icon" gorm:"-"`            // 菜单图标
	ShowMenu        bool   `json:"showMenu" gorm:"-"`        // 是否显示菜单
	HideMenu        bool   `json:"hideMenu" gorm:"-"`        // 是否隐藏菜单
	IgnoreKeepAlive bool   `json:"ignoreKeepAlive" gorm:"-"` // 是否忽略缓存
}

// 根据ID查询菜单，并预加载关联的角色信息
func GetMenuById(id int) (*Menu, error) {
	var dbMenu Menu
	err := DB.Where("id = ?", id).Preload("Roles").First(&dbMenu).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("菜单不存在")
		} else {
			return nil, fmt.Errorf("数据库错误：%w", err)
		}
	}
	return &dbMenu, nil
}

// UpdateOne 更新菜单
func (obj *Menu) UpdateOne() error {
	return DB.Updates(obj).Error
}

// CreateOne 创建菜单
func (obj *Menu) CreateOne() error {
	return DB.Create(obj).Error
}

// 获取菜单列表，包括其所有关联
func GetMenuAll() (menus []*Menu, err error) {
	err = DB.Find(&menus).Error
	return

}

// DeleteOne 删除菜单，包括其所有关联
func (obj *Menu) DeleteOne() error {
	return DB.Select(clause.Associations).Unscoped().Delete(obj).Error
}
