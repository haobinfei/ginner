package tools

import (
	"encoding/json"
	"fmt"
)

func JsonI2Struct(str, obj interface{}) {
	jsonStr := str.(string)
	Json2Struct(jsonStr, obj)
}

func Json2Struct(str string, obj interface{}) {
	// 将json转为结构体
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		panic(fmt.Sprintf("[Json2Struct]转换异常: %v", err))
	}
}
