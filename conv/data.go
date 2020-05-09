package conv

import (
	"encoding/json"
)

// ToJson 转换为JSON数据
func ToJson(data interface{}) []byte {
	buf, _ := json.Marshal(data)
	return buf
}

// StringSlice convert []interface to []string
func StringSlice(v []interface{}) []string {
	var result []string
	for _, t := range v {
		result = append(result, String(t))
	}
	return result
}

func Slice(v interface{}) []interface{} {
	switch res := v.(type) {
	case []interface{}:
		return res
	}
	return nil
}
