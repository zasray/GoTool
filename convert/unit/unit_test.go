package unit

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

// 测试存储单位转换器
func TestStorage(t *testing.T) {
	convertStorageUnit := ConvertStorageUnit{}
	s := fmt.Sprintf("%.0f", 1.899*1024*1024*1024)
	ret, _ := strconv.Atoi(s)
	cRet := convertStorageUnit.StringToInt(DEFAULT, MB, "1.899P")
	fmt.Println(fmt.Sprintf("预期结果：%d,实际返回：%d", ret, cRet))
	if math.Abs(float64(cRet-ret)) > 1 {
		t.Errorf("存储单位-TB-MB-转换错误")
	} else {
		t.Logf("存储单位-TB-MB-转换成功")
	}
}
