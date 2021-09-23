//go:build !windows
// +build !windows

package disk

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os/exec"
)

type LsblkRet struct {
	Blockdevices []Blockdevices `json:"blockdevices"`
}
type Children struct {
	Name       string     `json:"name"`
	Kname      string     `json:"kname"`
	Path       string     `json:"path"`
	Mountpoint string     `json:"mountpoint"`
	Label      string     `json:"label"`
	UUID       string     `json:"uuid"`
	Model      string     `json:"model"`
	Serial     string     `json:"serial"`
	Size       string     `json:"size"`
	Children   []Children `json:"children,omitempty"`
}
type Blockdevices struct {
	Name       string     `json:"name"`
	Kname      string     `json:"kname"`
	Path       string     `json:"path"`
	Mountpoint string     `json:"mountpoint"`
	Label      string     `json:"label"`
	UUID       string     `json:"uuid"`
	Model      string     `json:"model"`
	Serial     string     `json:"serial"`
	Size       string     `json:"size"`
	Children   []Children `json:"children,omitempty"`
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
	lsblkRet := LsblkRet{}
	err = json.Unmarshal(opBytes, &lsblkRet)
	ret := make([]DiskInfo, 0)
	for _, Blockdevices := range lsblkRet.Blockdevices {
		diskInfo := DiskInfo{}
		diskInfo.Model = Blockdevices.Model
		diskInfo.SerialNumber = Blockdevices.Serial
		diskInfo.Size = 0
		diskInfo.Children = make([]DiskChildren, 0)
		ret = append(ret, diskInfo)
	}
	return ret
}
