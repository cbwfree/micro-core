package mgo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

// MongoDB集合
type Table struct {
	name  string             // 表名称
	model reflect.Type       // 结构体反射类型
	index []mongo.IndexModel // 索引
	data  []interface{}      // 初始化数据
}

func (dt *Table) Name() string {
	return dt.name
}

func (dt *Table) Index() []mongo.IndexModel {
	return dt.index
}

func (dt *Table) Data() []interface{} {
	return dt.data
}

func (dt *Table) Model() reflect.Type {
	return dt.model
}

// AddIndex 设置索引
func (dt *Table) SetIndex(index []mongo.IndexModel) {
	dt.index = index
}

// AddIndex 追加索引
func (dt *Table) AddIndex(index []mongo.IndexModel) {
	dt.index = append(dt.index, index...)
}

func (dt *Table) SetData(data []interface{}) {
	dt.data = data
}

func NewTable(name string, model interface{}) *Table {
	vo := reflect.ValueOf(model)
	mt := &Table{
		name:  name,
		model: vo.Type(),
	}
	return mt
}
