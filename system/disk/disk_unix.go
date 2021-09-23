//go:build !windows
// +build !windows

package disk

import (
	"encoding/json"
	"github.com/zasray/GoTool/convert/unit"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

type LsblkRet struct {
	Blockdevices []Blockdevices `json:"blockdevices"`
}
type Blockdevices struct {
	Name       string         `json:"name"`
	Kname      string         `json:"kname"`
	Path       string         `json:"path"`
	Mountpoint string         `json:"mountpoint"`
	Label      string         `json:"label"`
	UUID       string         `json:"uuid"`
	Model      string         `json:"model"`
	Serial     string         `json:"serial"`
	Size       string         `json:"size"`
	Rota       string         `json:"rota"` //是否可旋转，判断ssd和hdd的关键指标
	Children   []Blockdevices `json:"children,omitempty"`
}

type DiskToolImpl struct {
}

//GetDiskList 获取硬盘列表,思路，执行
func (e *DiskToolImpl) GetDiskList() []DiskInfo {
	cmd := exec.Command("lsblk", "-JO")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	convertStorageUnit := unit.ConvertStorageUnit{}
	lsblkRet := LsblkRet{}
	err = json.Unmarshal(opBytes, &lsblkRet)
	ret := make([]DiskInfo, 0)
	for _, blockDevice := range lsblkRet.Blockdevices {
		diskInfo := DiskInfo{}
		diskInfo.Model = blockDevice.Model
		diskInfo.SerialNumber = blockDevice.Serial
		diskInfo.Size = float64(convertStorageUnit.StringToInt(unit.DEFAULT, unit.MB, blockDevice.Size))
		diskInfo.SSD = blockDevice.Rota == "0"
		diskInfo.System = false
		diskInfo.Children = make([]DiskChildren, 0)
		if len(blockDevice.Children) > 0 {
			for _, childPath := range blockDevice.Children {
				diskChildren := DiskChildren{}
				diskChildren.Path = childPath.Mountpoint
				if strings.Contains(strings.ToLower(diskChildren.Path), "/boot") {
					diskInfo.System = true
				}
				diskChildren.Size = float64(convertStorageUnit.StringToInt(unit.DEFAULT, unit.MB, childPath.Size))
				diskChildren.Free = 0
				diskInfo.Children = append(diskInfo.Children, diskChildren)
			}
		}
		ret = append(ret, diskInfo)
	}
	return ret
}
