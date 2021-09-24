//go:build windows
// +build windows

package disk

import (
	"fmt"
	"github.com/StackExchange/wmi"
	"regexp"
	"strings"
)

type Win32DiskDrive struct {
	Model        string
	Manufacturer string
	Signature    int
	TotalHeads   int
	SerialNumber string
	// 转换到G
	Size  int
	Index int
}

type Win32LogicalDiskToPartition struct {
	Antecedent *string
	Dependent  *string
}

type StorageInfo struct {
	Name       string
	Size       uint64
	FreeSpace  uint64
	FileSystem string
}

type DiskToolImpl struct {
}

//GetDiskList 获取硬盘列表
func (e *DiskToolImpl) GetDiskList() []DiskInfo {
	var win32DiskDrives []Win32DiskDrive
	var win32LogicalDiskToPartitions []Win32LogicalDiskToPartition
	var storageInfos []StorageInfo
	err := wmi.Query("select * from Win32_DiskDrive", &win32DiskDrives)
	if err != nil {
		panic(err)
	}
	err = wmi.Query("select * from Win32_LogicalDiskToPartition", &win32LogicalDiskToPartitions)
	if err != nil {
		panic(err)
	}

	err = wmi.Query("Select * from Win32_LogicalDisk", &storageInfos)
	if err != nil {
		panic(err)
	}

	// 开始构建
	ret := make([]DiskInfo, 0)
	for _, win32DiskDrive := range win32DiskDrives {
		diskInfo := DiskInfo{
			Model:        win32DiskDrive.Model,
			SerialNumber: win32DiskDrive.SerialNumber,
			Size:         float64(win32DiskDrive.Size / 1024.0 / 1024.0), //MB
			Children:     make([]DiskChildren, 0),
			SSD:          strings.Contains(strings.ToLower(win32DiskDrive.Model), strings.ToLower("NVMe")),
			System:       false,
		}
		//寻找分区信息
		for _, win32LogicalDiskToPartition := range win32LogicalDiskToPartitions {
			if !strings.Contains(*win32LogicalDiskToPartition.Antecedent, fmt.Sprintf("#%d,", win32DiskDrive.Index)) {
				continue
			}
			reg := regexp.MustCompile("DeviceID=\"(.*?)\"")
			matched := reg.MatchString(*win32LogicalDiskToPartition.Dependent)
			if matched {
				matchSubString := reg.FindStringSubmatch(*win32LogicalDiskToPartition.Dependent)
				size := float64(0)
				usedSpace := float64(0)
				freeSpace := float64(0)
				// 获取分区子盘的数据
				for _, storageInfo := range storageInfos {
					if storageInfo.Name == matchSubString[1] {
						if strings.Contains(strings.ToLower(storageInfo.Name), "c:") {
							// 系统盘
							diskInfo.System = true
						}
						size = float64(storageInfo.Size / 1024 / 1024)
						freeSpace = float64(storageInfo.FreeSpace / 1024 / 1024)
						usedSpace = size - freeSpace
					}
				}
				diskInfo.Children = append(diskInfo.Children, DiskChildren{
					Path: matchSubString[1],
					Size: size,
					Used: usedSpace,
					Free: freeSpace,
				})
			}
		}
		ret = append(ret, diskInfo)
	}
	return ret
}
