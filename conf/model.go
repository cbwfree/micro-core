package conf

import (
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
func convert(rows []*Model) map[string]interface{} {
	// 整理数据
	var conf = make(map[string]interface{})
	for _, row := range rows {
		switch row.Type {
		case "int":
			conf[row.Field] = conv.Int(row.Value)
		case "int32":
			conf[row.Field] = conv.Int32(row.Value)
		case "int64":
			conf[row.Field] = conv.Int64(row.Value)
		case "uint":
			conf[row.Field] = conv.Uint(row.Value)
		case "uint32":
			conf[row.Field] = conv.Uint32(row.Value)
		case "uint64":
			conf[row.Field] = conv.Uint64(row.Value)
		case "float32":
			conf[row.Field] = conv.Float32(row.Value)
		case "float64":
			conf[row.Field] = conv.Float64(row.Value)
		case "bool":
			conf[row.Field] = conv.Bool(row.Value)
		case "[]string":
			conf[row.Field] = strings.Split(row.Value, ",")
		case "[]int":
			var value []int
			for _, s := range strings.Split(row.Value, ",") {
				value = append(value, conv.Int(s))
			}
			conf[row.Field] = value
		case "[]int32":
			var value []int32
			for _, s := range strings.Split(row.Value, ",") {
				value = append(value, conv.Int32(s))
			}
			conf[row.Field] = value
		case "[]int64":
			var value []int64
			for _, s := range strings.Split(row.Value, ",") {
				value = append(value, conv.Int64(s))
			}
			conf[row.Field] = value
		default:
			conf[row.Field] = row.Value
		}
	}
	return conf
}
