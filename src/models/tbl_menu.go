package models

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	Name       string    `json:"name" gorm:"type:varchar(100);uniqueIndex;comment:名称"`
	Title      string    `json:"title" gorm:"comment:名称"`
	ParentMenu string    `json:"parentMenu" gorm:"varchar(5);comment:父级的id"`
	Meta       *MenuMeta `json:"meta" gorm:"-"`
	Icon       string    `json:"icon" gorm:"comment:图标"`
	//DbId       int64     `json:"dbId" gorm:"comment:-"` // 等同于数据库ID
	//Id         string  `json:"id" gorm:"type:varchar(100);uniqueIndex;comment:用于树拼接的字符串ID"`
	Type      string  `json:"type" gorm:"type:varchar(5);comment:类型 0=目录 1=子菜单"`
	Show      string  `json:"show" gorm:"type:varchar(5);comment:类型 0=禁用 1=启用"`
	OrderNo   int     `json:"orderNo" gorm:"comment:排序"`
	Component string  `json:"component" gorm:"type:varchar(50);comment:前端组件 菜单就是LAYOUT"`
	Redirect  string  `json:"redirect" gorm:"type:varchar(50);comment:显示路径"`
	Path      string  `json:"path" gorm:"type:varchar(50);comment:路由路径"`
	Remark    string  `json:"remark" gorm:"comment:用户描述"`
	HomePath  string  `json:"homePath" gorm:"comment:登陆后的默认首页"`
	Status    string  `json:"status" gorm:"default:1;comment:是否启用 0禁用 1启用"`
	Roles     []*Role `json:"roles" gorm:"many2many:role_menus;"`
}

type MenuMeta struct {
	Title string `json:"title" gorm:"-"` // 菜单标题
	Icon  string `json:"icon" gorm:"-"`  // 菜单图标
}

func (r *Menu) TableName() string {
	return "menus"
}

func (r *MenuMeta) TableName() string {
	return "Menus"
}
