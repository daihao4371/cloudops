package models

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Menu 结构体定义了菜单的数据库表字段及其关系
type Menu struct {
	Model        // 不用每次写ID 和 createAt了
	Name  string `json:"name" gorm:"type:varchar(100);uniqueIndex;comment:名称"`
	Title string `json:"title" gorm:"comment:名称" `
	//Title      string    `json:"title" gorm:"comment:名称" validate:"required,min=10,max=20"` // 测试一下 validate 和gorm在一起
	Pid        int       `json:"pId" gorm:"comment:父级的id"`
	ParentMenu string    `json:"parentMenu" gorm:"varchar(5);comment:父级的id"`
	Meta       *MenuMeta `json:"meta" gorm:"-"`
	Icon       string    `json:"icon" gorm:"comment:图标"`
	//DbId       uint      `json:"dbId" gorm:"-"` // 等同于数据库Id
	//Id        string  `json:"id" gorm:"-"`
	Type      string  `json:"type" gorm:"type:varchar(5);comment:类型 0=目录 1=子菜单"`
	Show      string  `json:"show" gorm:"type:varchar(5);comment:类型 0=禁用 1=启用"`
	OrderNo   int     `json:"orderNo" gorm:"comment:排序"`
	Component string  `json:"component" gorm:"type:varchar(50);comment:前端组件 菜单就是LAYOUT"`
	Redirect  string  `json:"redirect" gorm:"type:varchar(50);comment:显示路径"`
	Path      string  `json:"path" gorm:"type:varchar(50);comment:路由路径"`
	Remark    string  `json:"remark" gorm:"comment:用户描述"`
	HomePath  string  `json:"homePath" gorm:"comment:登陆后的默认首页"`
	Status    string  `json:"status" gorm:"default:1;comment:是否启用 0禁用 1启用"` //用户是否被冻结 1正常 2冻结
	Roles     []*Role `json:"roles" gorm:"many2many:role_menus;"`           // 多对多
	Children  []*Menu `json:"children" gorm:"-"`                            // 返回给前端的
	Key       uint    `json:"value"  gorm:"-"`
	Value     uint    `json:"key"  gorm:"-"`
}

// MenuMeta 结构体定义了菜单的元数据字段
type MenuMeta struct {
	Title           string `json:"title" gorm:"-"`
	Icon            string `json:"icon" gorm:"-"`
	ShowMenu        bool   `json:"showMenu" gorm:"-"`
	HideMenu        bool   `json:"hideMenu" gorm:"-"` // hideMenu=true 不显示 而不是showMenu
	IgnoreKeepAlive bool   `json:"ignoreKeepAlive" gorm:"-"`
}

// 根据ID查询菜单，并预加载关联的角色信息
func GetMenuById(id int) (*Menu, error) {
	var dbMenu Menu
	err := DB.Where("id = ?", id).Preload("Roles").First(&dbMenu).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("菜单不存在")
		}
		return nil, fmt.Errorf("数据库错误:%w", err)
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
