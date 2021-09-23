package disk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestFindDiskInfo(t *testing.T) {
	diskTool := DiskToolImpl{}
	diskList := diskTool.GetDiskList()
	if len(diskList) == 0 {
		t.Errorf("获取硬盘失败，硬盘数为%d，测试不通过！", len(diskList))
	} else {
		fmt.Println("成功了")
		fmt.Printf("%d個硬盤：%v\n", len(diskList), diskList)
		s, _ := json.Marshal(diskList)
		var out bytes.Buffer
		json.Indent(&out, s, "", "    ")
		fmt.Println(out.String())
		t.Logf("获取硬盘成功，硬盘数为%d，Json：%s", len(diskList), out.String())
	}
}
