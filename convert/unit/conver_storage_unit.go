package unit

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	KB = "K"
	MB = "M"
	GB = "G"
	TB = "T"
	PB = "P"
	EB = "E"
)

type ConvertStorageUnit struct {
}

func (e *ConvertStorageUnit) StringToInt(sourceUnit string, targetUnit string, data string) int {
	reg, _ := regexp.Compile("\\d+\\.?\\d*")
	value := float64(0)
	if reg.MatchString(data) {
		value, _ = strconv.ParseFloat(reg.FindStringSubmatch(data)[0], 64)
		if value > 0 {
			if sourceUnit == "" || sourceUnit != targetUnit {
				// 只有两个单位不统一时才需要：自行判断
				if sourceUnit == "" {
					//提取单位
					unitPatter, _ := regexp.Compile("(K|M|G|T|P|E)I?B?")
					if unitPatter.MatchString(strings.ToUpper(data)) {
						sourceUnit = unitPatter.FindStringSubmatch(data)[1]
					}
				}

				//如果还提取不了原单位，则直接报错
				if sourceUnit == "" {
					panic(fmt.Sprintf("can not get sourceUnit by:%s", data))
				} else {
					switch sourceUnit {
					case KB:
						value = value / 1024
					case GB:
						value = value * 1024
					case TB:
						value = value * 1024 * 1024
					case PB:
						value = value * 1024 * 1024 * 1024
					case EB:
						value = value * 1024 * 1024 * 1024 * 1024
					default:
						// 默认就是MB 不需要转换
						value = value * 1
					}

					if targetUnit == "" {
						targetUnit = "MB"
					}
					switch targetUnit {
					case KB:
						value = value * 1024
					case GB:
						value = value / 1024
					case TB:
						value = value / 1024 / 1024
					case PB:
						value = value / 1024 / 1024 / 1024
					case EB:
						value = value / 1024 / 1024 / 1024 / 1024
					}
				}
			}
		}
	}
	return int(value)
}
