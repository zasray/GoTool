package unit

const (
	DEFAULT = ""
)

//ConvertUnit 转换工具类
type ConvertUnit interface {
	//StringToInt 字符串转数字
	StringToInt(sourceUnit string, targetUnit string, data string) int
}
