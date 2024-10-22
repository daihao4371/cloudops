package models

import (
	"cloudops/src/common"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 用户相关的数据库字段
type User struct {
	Model               // 不用每次写ID 和 createAt了
	UserId       int    `json:"userId" gorm:"comment:用户id"`
	Username     string `json:"username" gorm:"type:varchar(100);uniqueIndex;comment:用户登录名"` // 用户登录名 uniqueIndex 代表唯一索引
	Password     string `json:"password"  gorm:"comment:用户登录密码"`
	RealName     string `json:"realName" gorm:"type:varchar(100);uniqueIndex;comment:用户昵称 中文名"`
	Desc         string `json:"desc" gorm:"comment:用户描述"`
	Mobile       string `json:"mobile" gorm:"comment:手机号"`
	FeiShuUserId string `json:"feiShuUserId" gorm:"comment:飞书userid  获取方式请看 https://open.feishu.cn/document/home/user-identity-introduction/open-id"`
	// 添加一些字段 用来支持服务账号
	AccountType         int          `json:"accountType" gorm:"default:1;comment:用户是否是服务账号 1普通用户 2服务账号"`
	ServiceAccountToken string       `json:"serviceAccountToken"`
	HomePath            string       `json:"homePath" gorm:"comment:登陆后的默认首页"`
	Enable              int          `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"` //用户是否被冻结 1正常 2冻结
	Roles               []*Role      `json:"roles" gorm:"many2many:user_roles;"`              // 多对多
	OpsNodes            []*StreeNode `json:"ops_nodes" gorm:"many2many:ops_admins;"`          // 多对多
	RdAdminNodes        []*StreeNode `json:"rd_admin_nodes" gorm:"many2many:rd_admins;"`      // 多对多
	RdMemberNodes       []*StreeNode `json:"rd_member_nodes" gorm:"many2many:rd_members;"`    // 多对多
	RolesFront          []string     `json:"rolesFront" gorm:"-"`                             // 给前端用的中间字段
}

// CheckUserPassword 验证用户登录密码
func CheckUserPassword(ru *UserLoginRequest) (*User, error) {
	var dbUser User
	// 根据用户名查询用户信息，并预加载用户角色
	err := DB.Where("username = ?", ru.Username).Preload("Roles").First(&dbUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户名不存在")
		}
		return nil, fmt.Errorf("数据库错误:%w", err)
	}

	// 对比用户密码
	ok := common.BcryptCheck(ru.Password, dbUser.Password)
	if ok {

		return &dbUser, nil
	}
	return nil, fmt.Errorf("密码错误")
}

// GetUserByUserName 根据用户名获取用户信息
func GetUserByUserName(userName string) (*User, error) {
	var dbUser User
	// 根据用户名查询用户信息，并预加载用户角色及其菜单
	err := DB.Where("username = ?", userName).Preload("Roles.Menus").First(&dbUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户名不存在")
		}
		return nil, fmt.Errorf("数据库错误:%w", err)
	}
	return &dbUser, nil
}

// GetUserByRealName 根据用户花名名查询用户
func GetUserByRealName(realName string) (*User, error) {
	var dbUser User
	err := DB.Where("real_name = ?", realName).Preload("Roles").First(&dbUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("数据库错误: %w", err)
	}
	return &dbUser, nil
}

// GetUserById 根据用户ID获取用户信息
func GetUserById(id int) (*User, error) {

	var dbUser User
	err := DB.Where("id = ?", id).Preload("Roles").First(&dbUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("数据库错误:%w", err)
	}
	return &dbUser, nil

}

// UpdateOne 更新User对象的本体及其关联的角色
func (obj *User) UpdateOne(roles []*Role) error {
	err1 := DB.Where("id = ?", obj.ID).Updates(obj).Error
	err2 := DB.Model(obj).Association("Roles").Replace(roles)
	if err1 == nil && err2 == nil {
		return nil
	} else {
		return fmt.Errorf("更新本体:%w 更新关联:%w", err1, err2)
	}

}

// FirstOrCreate 方法用于查询数据库中是否存在具有相同Username属性的User对象，
// 如果存在则更新obj对象的属性值，否则在数据库中创建新的User对象。
// 方法返回值为error类型，如果操作成功则返回nil，否则返回具体的错误信息。
func (obj *User) FirstOrCreate() error {
	return DB.Debug().Where(User{Username: obj.Username}).FirstOrCreate(obj).Error
}

// UpdateRoles 更新用户关联的角色
func (obj *User) UpdateRoles(roles []*Role) error {
	return DB.Model(obj).Association("Roles").Replace(roles)

}

// CreateOne 创建用户，并返回错误信息。
func (obj *User) CreateOne() error {
	return DB.Create(obj).Error
}

// // GetUserAll 获取数据库中获取所有用户及其关联的角色，并返回用户切片和错误信息。
func GetUserAll() (users []*User, err error) {
	err = DB.Preload("Roles").Find(&users).Error
	return
}

// GetUserByName 根据用户名获取用户信息
func GetUserByName(name string) (*User, error) {
	var dbUser User
	err := DB.Where("username = ?", name).Preload("Roles").First(&dbUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("数据库错误:%w", err)
	}
	return &dbUser, nil
}

// 删除用户
func (obj *User) DeleteOne() error {
	return DB.Select(clause.Associations).Unscoped().Delete(obj).Error
}
