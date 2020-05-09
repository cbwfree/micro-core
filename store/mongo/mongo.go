// MongoDB 连接
// * 查找单个文档时, 如果未找到文件, 则会返回 ErrNoDocuments 错误
// * 查找多个文档时, 如果未找到任何文档, 则会返回 ErrNilDocument 错误
// * bson.M 是无序的 doc 描述
// * bson.D 是有序的 doc 描述
// * bsonx.Doc 是类型安全的 doc 描述
package mgo

import (
	"context"
	"errors"
	"github.com/micro/go-micro/v2/util/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"strings"
)

// MongoDB 数据存储
type Store struct {
	dbname string
	opts   *Options
	client *mongo.Client
}

func (ms *Store) With(opts ...Option) {
	ms.opts.With(opts...)
}

func (ms *Store) DbName() string {
	return ms.dbname
}

func (ms *Store) Opts() *Options {
	return ms.opts
}

func (ms *Store) Connect() error {
	if ms.client != nil {
		return nil
	}

	if ms.opts.Uri == "" {
		ms.opts.RawUrl = "mongodb://0.0.0.0:27017,0.0.0.0:27018,0.0.0.0:27019/?replicaSet=rs1"
	} else if !strings.HasPrefix(ms.opts.Uri, "mongodb://") {
		ms.opts.RawUrl = "mongodb://" + ms.opts.Uri
	} else {
		ms.opts.RawUrl = ms.opts.Uri
	}

	opts := options.Client().
		SetMinPoolSize(ms.opts.MinPoolSize).
		SetMaxPoolSize(ms.opts.MaxPoolSize).
		SetConnectTimeout(ms.opts.ConnectTimeout).
		SetSocketTimeout(ms.opts.SocketTimeout).
		SetMaxConnIdleTime(ms.opts.MaxConnIdleTime).
		SetRetryWrites(true).
		SetRetryReads(true).
		ApplyURI(ms.opts.RawUrl)
	if *opts.ReplicaSet == "" {
		return errors.New("this system only supports replica sets. example: mongodb://0.0.0.0:27017,0.0.0.0:27018,0.0.0.0:27019/?replicaSet=rs1")
	}

	if mc, err := mongo.Connect(nil, opts); err != nil {
		return err
	} else {
		ms.client = mc
	}

	// 检查MongoDB连接
	if err := ms.client.Ping(nil, readpref.Primary()); err != nil {
		return err
	}

	log.Debugf("Store [mongodb] Connect to %s", ms.opts.RawUrl)

	return nil
}

// 关闭连接
func (ms *Store) Disconnect() error {
	if ms.client == nil {
		return nil
	}

	if err := ms.client.Disconnect(nil); err != nil {
		return err
	}

	ms.client = nil

	return nil
}

// Client 获取客户端
func (ms *Store) Client() *mongo.Client {
	return ms.client
}

// Database 获取数据库对象
func (ms *Store) D(dbname ...string) *mongo.Database {
	if len(dbname) > 0 && dbname[0] != "" {
		return ms.client.Database(dbname[0])
	}
	return ms.client.Database(ms.dbname)
}

// Collection 获取集合对象
func (ms *Store) C(name string, dbname ...string) *mongo.Collection {
	if len(dbname) > 0 && dbname[0] != "" {
		return ms.client.Database(dbname[0]).Collection(name)
	}
	return ms.client.Database(ms.dbname).Collection(name)
}

// CloneCollection 克隆集合对象
func (ms *Store) CloneC(name string, dbname ...string) (*mongo.Collection, error) {
	return ms.C(name, dbname...).Clone()
}

// 获取自增ID
func (ms *Store) GetIncId(id string) (int64, error) {
	return GetIncId(context.Background(), ms.D(), id)
}

// 获取集合列表
func (ms *Store) ListCollectionNames(dbname ...string) ([]string, error) {
	return ms.D(dbname...).ListCollectionNames(context.Background(), bson.M{})
}

// 分段获取数据
func (ms *Store) Scan(dbName, tabName string, cur, size int64, filter interface{}, result interface{}, fn ...func(opts *options.FindOptions) *options.FindOptions) *Scan {
	var scan *Scan
	var closure = func(sctx mongo.SessionContext) error {
		col := sctx.Client().Database(dbName).Collection(tabName)

		count, _ := col.CountDocuments(sctx, filter)
		scan = NewScan(count, cur, size)

		if count > 0 {
			opts := scan.FindOptions()
			if len(fn) > 0 {
				opts = fn[0](opts)
			}
			cur, err := col.Find(sctx, filter, opts)
			if err != nil {
				return err
			}
			if err := cur.All(nil, result); err != nil {
				return err
			}
		}

		return nil
	}
	if err := ms.Client().UseSession(context.Background(), closure); err != nil {
		return new(Scan)
	}
	return scan
}

// 实例化MongoDB存储
func NewStore(dbname string, opts ...Option) *Store {
	ms := &Store{
		dbname: dbname,
		opts:   newOptions(opts...),
	}
	return ms
}
