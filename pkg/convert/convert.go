package convert

import (
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// 字符串转数字
func StrToNum(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// JSON转字符串
func Marshal[T any](j T) string {
	a, _ := json.Marshal(j)
	return string(a)
}

// 字符串转JSON
func Unmarshal[T any](s string, data T) {
	_ = json.Unmarshal([]byte(s), &data)
}
