package disk

// DiskInfo 硬盘信息
type DiskInfo struct {
	// 型号
	Model string `json:"model"`
	// 是否是ssd
	SSD bool `json:"ssd"`
	// 是否是系统盘，根据当前的理解，windows下含有C盘的硬盘是系统盘、linux下含有/boot的硬盘是系统盘
	System bool `json:"system"`
	// 序列号(硬盘的唯一判断标准）
	SerialNumber string `json:"serial_number"`
	// 大小 MB
	Size float64 `json:"size"`
	// 分区（路径）
	Children []DiskChildren `json:"children"`
}

// DiskChildren 硬盘子目录信息
type DiskChildren struct {
	// 路径
	Path string `json:"path"`
	// 大小 MB
	Size float64 `json:"size"`
	// 已用空间
	Used float64 `json:"used"`
	// 剩余空间
	Free float64 `json:"free"`
}

// 硬盘工具
type DiskTool interface {
	GetDiskList() []DiskInfo
}
