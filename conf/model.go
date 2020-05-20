package conf

import (
	"encoding/json"
	"github.com/cbwfree/micro-core/conv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"strings"
)

// 系统配置数据模型
type Model struct {
	Field string `bson:"field"` // 设置字段
	Type  string `bson:"type"`  // 数据类型
	Value string `bson:"value"` // 设置值
}

var (
	// 设置索引
	Indexes = []mongo.IndexModel{
		{
			Keys: bsonx.Doc{
				{Key: "field", Value: bsonx.Int32(1)},
			},
			Options: options.Index().SetUnique(true),
		},
	}
)

// 转换模型数据为对应类型数据
func Convert(t string, v string) interface{} {
	switch t {
	case "bool":
		return conv.Bool(v)
	case "int":
		return conv.Int(v)
	case "int32":
		return conv.Int32(v)
	case "int64":
		return conv.Int64(v)
	case "float32":
		return conv.Float32(v)
	case "float64":
		return conv.Float64(v)
	case "[]string":
		return strings.Split(v, ",")
	case "[]int":
		var value []int
		for _, s := range strings.Split(v, ",") {
			value = append(value, conv.Int(s))
		}
		return value
	case "[]int32":
		var value []int32
		for _, s := range strings.Split(v, ",") {
			value = append(value, conv.Int32(s))
		}
		return value
	case "[]int64":
		var value []int64
		for _, s := range strings.Split(v, ",") {
			value = append(value, conv.Int64(s))
		}
		return value
	default:
		return conv.String(v)
	}
}

func toDataJson(rows []*Model) []byte {
	var data = make(map[string]interface{})
	for _, row := range rows {
		data[row.Field] = Convert(row.Type, row.Value)
	}
	b, _ := json.Marshal(data)
	return b
}
