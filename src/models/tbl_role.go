package models

import (
	"fmt"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	OrderNo   int     `json:"orderNo" gorm:"comment:排序"`
	RoleName  string  `json:"roleName" gorm:"type:varchar(100);uniqueIndex;comment:角色中文名称"`
	RoleValue string  `json:"roleValue"  gorm:"type:varchar(100);uniqueIndex;comment:角色值"`
	Remark    string  `json:"remark" gorm:"comment:用户描述"`
	HomePath  string  `json:"homePath" gorm:"comment:登陆后的默认首页"`
	Status    string  `json:"status" gorm:"default:1;comment:角色是否被冻结 1正常 2冻结"`
	Users     []*User `json:"users" gorm:"many2many:user_roles;"`
	Menus     []*Menu `json:"menus" gorm:"many2many:role_menus;"` // 多对多关系，角色对应的菜单
	Apis      []*Api  `json:"apis" gorm:"many2many:role_apis;"`
	MenuIds   []int   `json:"menuIds" gorm:"-"`
	ApiIds    []int   `json:"apiIds" gorm:"-"`
}

// 获取所有角色,用户和API
func GetRoleAll() (roles []*Role, err error) {
	err = DB.Preload("Menus").Preload("Users").Preload("Apis").Find(&roles).Error
	return
}

// 根据角色获取角色信息，
func GetRoleByRoleValue(roleValue string) (*Role, error) {
	var dbRole Role
	err := DB.Where("role_value = ?", roleValue).Preload("Menus").First(&dbRole).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("角色不存在")
		}
		return nil, fmt.Errorf("数据库错误:%w", err)
	}
	return &dbRole, nil
}

// 根据角色ID获取角色信息
func GetRoleById(id int) (*Role, error) {
	var dbrole Role
	err := DB.Where("id = ?", id).Preload("Menus").First(&dbrole).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("角色不存在")
		}
		return nil, fmt.Errorf("数据库错误:%w", err)
	}
	return &dbrole, err
}

// 创建一个新角色
func (obj *Role) CreateOne() error {
	return DB.Create(&obj).Error
}

// 更新角色的菜单信息
func (obj *Role) UpdateMenus(meuns []*Menu) error {
	// 更新角色信息
	err1 := DB.Where("id = ?", obj.ID).Updates(obj).Error
	// 更新角色菜单
	err2 := DB.Model(&obj).Association("Menus").Replace(meuns)
	if err1 == nil && err2 == nil {
		return nil
	} else if err1 != nil && err2 == nil {
		return fmt.Errorf("更新角色信息失败:%w", err1)
	} else if err1 == nil && err2 != nil {
		return fmt.Errorf("更新角色菜单失败:%w", err2)
	} else {
		return fmt.Errorf("更新角色信息失败:%w,更新角色菜单失败:%w", err1, err2)
	}

}
