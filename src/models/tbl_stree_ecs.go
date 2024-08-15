package models

// 字段说明 https://help.aliyun.com/zh/ecs/developer-reference/api-describeinstances?spm=a2c4g.11186623.0.i1#t9865.html
type ResourceEcs struct {
	Model // 不用每次写ID 和 createAt了
	ResourceCommon
	// 核心字段
	// uid 代表不会变的 id ，除非被删除

	// 常见字段
	OsType       string `json:"OsType" gorm:"comment:操作系统类型 win linux" `
	VmType       int    `json:"VmType" gorm:"default:1;comment:角色是否被冻结 =1云厂商虚拟设备 =2物理设备 agent上报的"`
	InstanceType string `json:"InstanceType" gorm:"comment:实例类型 ecs.g8a.2xlarge  https://www.alibabacloud.com/help/zh/ecs/user-guide/overview-of-instance-families " `

	// 资源
	Cpu               int         `json:"Cpu" gorm:"comment:vCPU数"`
	Memory            int         `json:"Memory" gorm:"comment:内存大小，单位为GiB。"`
	Disk              int         `json:"Disk" gorm:"comment:磁盘打下，单位为GiB。"`
	OSName            string      `json:"OSName" gorm:"comment:实例的操作系统名称。 CentOS 7.4 64 位"`
	ImageId           string      `json:"ImageId" gorm:"comment:镜像模板" `
	Hostname          string      `json:"Hostname" gorm:"type:varchar(100);comment:主机名" `
	NetworkInterfaces StringArray `json:"NetworkInterfaces" gorm:"comment: []弹性网卡id集合"`
	DiskIds           StringArray `json:"DiskIds" gorm:"comment: [d-bp67acfmxazb4p]云盘ID"`

	// 有关时间的
	StartTime       string `json:"StartTime" gorm:"comment:实例最近一次的启动时间。以ISO 8601为标准，并使用UTC+0时间，格式为yyyy-MM-ddTHH:mmZ。更多信息，请参见ISO 8601。"  `
	AutoReleaseTime string `json:"AutoReleaseTime" `
	LastInvokedTime string `json:"LastInvokedTime" `

	BindNodes []*StreeNode `json:"bind_nodes" gorm:"many2many:bind_ecss;"` // 多对多

}

type EcsBuyWorkOrder struct {
	Vendor         string `json:"vendor"`
	Num            int    `json:"num"`
	BindLeafNodeId int    `json:"bindLeafNodeId"`
	InstanceType   string `json:"instance_type"`
	Hostnames      string `json:"hostnames"` // split \n 多条记录
}
