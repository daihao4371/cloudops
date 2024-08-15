package models

type ResourceElb struct {
	Model // 不用每次写ID 和 createAt了
	ResourceCommon

	// 重点字段
	LoadBalancerType   string `json:"LoadBalancerType" gorm:"comment:类型有 nlb alb clb"`
	BandwidthCapacity  int    `json:"BandwidthCapacity" gorm:"comment:带宽包上限 50Mb 100Mb"`
	AddressType        string `json:"AddressType" gorm:"comment:公网类型还是内网类型"`
	DNSName            string `json:"DNSName" gorm:"comment: dns解析地址"`
	BandwidthPackageId string `json:"BandwidthPackageId" gorm:"comment:绑定的带宽包"`
	CrossZoneEnabled   bool   `json:"CrossZoneEnabled" `

	BindNodes []*StreeNode `json:"bind_nodes" gorm:"many2many:bind_elbs;"` // 多对多
}
