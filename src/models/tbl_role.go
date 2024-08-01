package models

import (
	"cloudops/src/config"
	"encoding/json"
	"fmt"
	"github.com/gammazero/workerpool"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

// 角色表
type Role struct {
	Model             // 不用每次写ID 和 createAt了
	OrderNo   int     `json:"orderNo" gorm:"comment:排序"`
	RoleName  string  `json:"roleName" gorm:"type:varchar(100);uniqueIndex;comment:角色中文名称"` // 用户登录名 index 代表索引
	RoleValue string  `json:"roleValue"  gorm:"type:varchar(100);uniqueIndex;comment:角色值"`
	Remark    string  `json:"remark" gorm:"comment:用户描述"`
	HomePath  string  `json:"homePath" gorm:"comment:登陆后的默认首页"`
	Status    string  `json:"status" gorm:"default:1;comment:角色是否被冻结 1正常 2冻结"` //用户是否被冻结 1正常 2冻结
	Users     []*User `json:"users" gorm:"many2many:user_roles;"`              // 多对多
	Menus     []*Menu `json:"menus" gorm:"many2many:role_menus;"`              // 多对多
	Apis      []*Api  `json:"apis" gorm:"many2many:role_apis;"`                // 多对多
	MenuIds   []int   `json:"menuIds" gorm:"-"`                                // 给前端用的
	ApiIds    []int   `json:"apiIds" gorm:"-"`                                 // 给前端用的
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
func GetRoleById(id uint) (*Role, error) {
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

// 更新角色的API信息
func (obj *Role) UpdateApis(apis []*Api, sc *config.ServerConfig) error {
	err1 := DB.Where("id = ?", obj.ID).Updates(obj).Error
	err2 := DB.Model(obj).Association("Apis").Replace(apis)
	rules := [][]string{}
	// 先获取所有的policy
	wp := workerpool.New(100)

	for _, api := range apis {
		api := api
		wp.Submit(func() {

			oneRule := []string{
				obj.RoleValue,
				api.Path,
				api.Method,
			}
			// 处理ALL的case
			if api.Method == "ALL" {
				methods := []string{
					"GET",
					"POST",
					"DELETE",
				}
				for _, m := range methods {

					// casbin先获取这个policy是否存在 = 判断这个3个值有没有权限
					// 如果pass=true 代表存在就不需要再添加了

					pass, err := CasbinCheckPermission(obj.RoleValue, api.Path, m)
					if pass {
						continue
					}
					_, err = CasbinAddOnePolicy(obj.RoleValue, api.Path, m)
					if err != nil {
						sc.Logger.Error("CasbinAddOnePolicy错误",
							zap.Error(err),
							zap.String("角色", obj.RoleValue),
							zap.String("api.Path", api.Path),
							zap.String("api.Method", api.Method),
						)
					}
				}
			} else {
				pass, err := CasbinCheckPermission(obj.RoleValue, api.Path, api.Method)
				if pass {
					return
				}
				_, err = CasbinAddOnePolicy(obj.RoleValue, api.Path, api.Method)
				if err != nil {
					sc.Logger.Error("CasbinAddOnePolicy错误",
						zap.Error(err),
						zap.String("角色", obj.RoleValue),
						zap.String("api.Path", api.Path),
						zap.String("api.Method", api.Method),
					)
				}
			}
			rules = append(rules, oneRule)
		})

	}
	wp.StopWait()

	if err1 == nil && err2 == nil {
		return nil
	} else {
		return fmt.Errorf("更新本体:%w 更新关联:%w", err1, err2)
	}
}

// DeleteOne 方法从数据库中删除Role实例obj。
func (obj *Role) DeleteOne() error {
	return DB.Delete(clause.Associations).Unscoped().Delete(obj).Error
}

/*
UnmarshalJSON 实现了json.Unmarshaler接口，用于反序列化Role结构体中的MenuIds字段
data: JSON格式数据
返回值：
error: 如果反序列化过程中发生错误，则返回相应的错误信息，否则返回nil
*/
func (obj *Role) UnmarshalJSON(data []byte) error {
	var temp struct {
		MenuIds []interface{} `json:"menuIds"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	for _, id := range temp.MenuIds {
		switch v := id.(type) {
		case float64:
			obj.MenuIds = append(obj.MenuIds, int(v))
		case string:
			menuId, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			obj.MenuIds = append(obj.MenuIds, menuId)
		default:
			return fmt.Errorf("unexpected type for MenuIds: %T", v)
		}
	}
	return nil
}
