package mgo

import (
	"context"
	"github.com/cbwfree/micro-core/fn"
	log "github.com/micro/go-micro/v2/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

// 数据库初始化
type Tables struct {
	sync.RWMutex
	tables map[string]*Table
}

// Tables 返回所有数据表
func (mts *Tables) Tables() map[string]*Table {
	return mts.tables
}

// Count 返回集合数量
func (mts *Tables) Count() int {
	mts.RLock()
	defer mts.RUnlock()

	return len(mts.tables)
}

// Get 获取指定集合
func (mts *Tables) Get(name string) *Table {
	mts.RLock()
	defer mts.RUnlock()

	for n, col := range mts.tables {
		if n == name {
			return col
		}
	}

	return nil
}

// Add 添加新集合
func (mts *Tables) Add(name string, model interface{}, index []mongo.IndexModel, data []interface{}) *Table {
	tab := NewTable(name, model)
	if len(index) > 0 {
		tab.SetIndex(index)
	}
	if len(data) > 0 {
		tab.SetData(data)
	}

	mts.Lock()
	mts.tables[name] = tab
	mts.Unlock()

	return tab
}

// 设置自增初始化数据
func (mts *Tables) SetAutoIdData(data []interface{}) {
	mts.Add(AutoIncIdName, AutoIncId{}, nil, data)
}

// Init 初始化数据库
func (mts *Tables) Check(mdb *Store) error {
	if mts.Count() == 0 {
		return nil
	}

	log.Debugf("[%s] check collections ...", mdb.DbName())

	// 获取集合列表
	names, err := mdb.ListCollectionNames(mdb.DbName())
	if err != nil {
		return err
	}

	var tables []string
	var closure = func(sctx mongo.SessionContext) error {
		cdb := mdb.Client().Database(mdb.DbName())

		for _, tab := range mts.tables {
			// 判断集合是否已存在
			if !fn.InStrSlice(tab.name, names) {
				// 获取MongoDB的集合对象
				col := cdb.Collection(tab.name)

				// 创建索引
				if len(tab.index) > 0 {
					if _, err := col.Indexes().CreateMany(sctx, tab.Index()); err != nil {
						log.Errorf("create table [ %s ] index failure, error: %s", tab.name, err.Error())
					} else {
						log.Infof("init table [ %s ] index success. total: %d", tab.name, len(tab.index))
					}
				}

				// 初始化数据
				if len(tab.data) > 0 {
					if _, err := col.InsertMany(sctx, tab.data); err != nil {
						log.Errorf("init table [ %s ] data failure, error: %s", tab.name, err.Error())
					} else {
						log.Infof("init table [ %s ] data success. total: %d", tab.name, len(tab.data))
					}
				}

				tables = append(tables, tab.name)
			}
		}

		return nil
	}
	if err := mdb.Client().UseSession(context.Background(), closure); err != nil {
		return err
	}

	if len(tables) > 0 {
		log.Infof("[%s] successfully initialize %d table ...", mdb.DbName(), len(tables))
	}

	return nil
}

// 实例化数据库模型
func NewTables() *Tables {
	d := &Tables{
		tables: make(map[string]*Table),
	}
	return d
}
