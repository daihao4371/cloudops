package models

type ResourceRds struct {
	Model // 不用每次写ID 和 createAt了
	ResourceCommon

	Engine            string `json:"Engine" gorm:"comment:类型有 mysql mariadb postgresql"`
	DBInstanceNetType string `json:"DBInstanceNetType" gorm:"comment:实例的网络连接类型，取值：Internet：外网连接 Intranet：内网连接"`
	DBInstanceClass   string `json:"DBInstanceClass" gorm:"comment:实例的实例规格，取值：rds.mys2.small https://help.aliyun.com/zh/rds/product-overview/primary-apsaradb-rds-instance-types"`
	DBInstanceType    string `json:"DBInstanceType" gorm:"comment:是否主备，Primary：主实例 Readonly：只读实例Guard：灾备实例Temp：临时实例"`
	EngineVersion     string `json:"EngineVersion" gorm:"comment:数据库版本。，取值：8.0 5.7"`
	MasterInstanceId  string `json:"MasterInstanceId" xml:"MasterInstanceId"`
	DBInstanceStatus  string `json:"DBInstanceStatus" xml:"DBInstanceStatus"`
	ReplicateId       string `json:"ReplicateId" xml:"ReplicateId"`

	BindNodes []*StreeNode `json:"bind_nodes" gorm:"many2many:bind_rdss;"` // 多对多
}
