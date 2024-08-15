package models

// 资源表的公共字段
type ResourceCommon struct {
	InstanceId    string `json:"InstanceId" gorm:"uniqueIndex;type:varchar(100);comment:实例ID  i-bp67acfmxazb4p****" `
	InstanceName  string `json:"title" gorm:"uniqueIndex;type:varchar(100);comment:实例名称，支持使用通配符*进行模糊搜索" `
	Hash          string `json:"hash" gorm:"uniqueIndex;type:varchar(200);comment:增量更新的"`
	Vendor        string `json:"Vendor" gorm:"comment: 云厂商 阿里云 华为云aws"` //用户是否被冻结 1正常 2冻结
	CreateByOrder bool   `json:"createByOrder" gorm:"comment:目的是工单创建的不要被增量更新给删掉"`
	VpcId         string `json:"VpcId" gorm:"comment:专有网络VPC ID。"`
	ZoneId        string `json:"ZoneId" gorm:"comment:实例所属可用区 cn-hangzhou-g 。"`
	Env           string `json:"Env" gorm:"comment:环境标识 dev=开发环境 press=压测环境 stage=预发环境 prod=生产环境"`
	PayType       string `json:"PayType" gorm:"comment:付费类型，取值：Postpaid：按量付费 Prepaid：包年包月"`

	//Region string     `json:"Region" gorm:""`
	Status      string      `json:"Status"  gorm:"comment:实例状态。取值范围： Pending：创建中。Running：运行中。Stopped：已停止。"`
	Description string      `json:"Description" gorm:"comment:实例描述。。 CentOS 7.4 64 位"`
	Tags        StringArray `json:"Tags"  gorm:"comment: [k1=v1,k2=v2]标签集合"`

	// 字符串数组类型的在这里
	SecurityGroupIds StringArray `json:"SecurityGroupIds"  gorm:"comment: [sg-bp67acfmxazb4p ]实例 安全组ID"`
	PrivateIpAddress StringArray `json:"PrivateIpAddress"   gorm:"comment: [1.1.1.1 ]私有IP地址"`
	PublicIpAddress  StringArray `json:"PublicIpAddress"   gorm:"comment: [1.1.1.1 ]公网IP地址"`
	IpAddr           string      `json:"ipAddr"  gorm:"comment: [1.1.1.1 ]公网IP地址"`
	CreationTime     string      `json:"CreationTime" gorm:"comment:  2017-12-10T04:04Z 实例或创建时间，注意和同步到数据库的CreateAt 2个概念。以ISO 8601为标准，并使用UTC+0时间，格式为yyyy-MM-ddTHH:mmZ。更多信息，请参见ISO 8601。"`

	// 绑定的叶子节点
	Key string `json:"key" gorm:"-"` //给前端用的
}
