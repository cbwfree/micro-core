package srv

import (
	"context"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/cbwfree/micro-core/conv"
	mgo "github.com/cbwfree/micro-core/store/mongo"
	rds "github.com/cbwfree/micro-core/store/redis"
	"github.com/cbwfree/micro-core/web"
	"github.com/go-redis/redis/v7"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/server"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func APP() *App {
	return app
}

func OPTS() *Options {
	return app.opts
}

func New(srvName string, flags ...[]cli.Flag) {
	APP().New(srvName, flags...)
}

func Run() error {
	return APP().Srv().Run()
}

func Close() {
	APP().Close()
}

func Init(opts ...micro.Option) {
	APP().Init(opts...)
}

func With(with ...WithAPP) {
	APP().With(with...)
}

// Srv 获取当前服务对象
func Srv() micro.Service {
	return APP().Srv()
}

// Id 获取当前服务ID
func Id() string {
	return APP().Id()
}

// Name 获取当前服务名称 (此名称包含集群标识)
func Name() string {
	return APP().Name()
}

// NameId 获取当前服务完整名称 (即服务名+节点ID)
func NameId() string {
	return APP().NameId()
}

// Server 获取服务器对象
func Server() server.Server {
	return APP().Srv().Server()
}

// Client  获取客户端对象
func Client() client.Client {
	return APP().Srv().Client()
}

// RS Redis存储
func RS() *rds.Store {
	return APP().Redis
}

// RSC Redis存储
func RSC() *redis.Client {
	return APP().Redis.Client()
}

// RSDo 执行Redis指令
func RSDo(args ...interface{}) *redis.Cmd {
	return APP().Redis.Do(args...)
}

// 分布式锁 (默认生存周期, 默认重试时间)
func RSLock(key string, index interface{}) (*redislock.Lock, error) {
	return RSLockBackoff(key, index, 3*time.Second, 100*time.Millisecond)
}

// 分布式锁 (指定生存周期)
func RSLockTTL(key string, index interface{}, ttl time.Duration) (*redislock.Lock, error) {
	return RSLockBackoff(key, index, ttl, 100*time.Millisecond)
}

// 分布式锁 (指定生存周期, 设置重试时间)
func RSLockBackoff(key string, index interface{}, ttl time.Duration, backoff time.Duration) (*redislock.Lock, error) {
	opt := &redislock.Options{}
	if backoff > 0 {
		opt.RetryStrategy = redislock.LinearBackoff(backoff)
	}
	keyName := fmt.Sprintf("LOCK:%s:%s", key, conv.String(index))
	return RS().Locker().Obtain(keyName, ttl, opt)
}

// MS MongoDB存储
func MS() *mgo.Store {
	return APP().Mongo
}

// MS MongoDB存储
func MSC() *mongo.Client {
	return APP().Mongo.Client()
}

// MSD MongoDB存储
func MSD(dbname ...string) *mongo.Database {
	return APP().Mongo.D(dbname...)
}

// MST MongoDB存储
func MST(table string, dbname ...string) *mongo.Collection {
	return APP().Mongo.C(table, dbname...)
}

// 发布广播
func PubCtx(ctx context.Context, name string, msg interface{}, opts ...client.PublishOption) error {
	return APP().PubCtx(ctx, name, msg, opts...)
}

// 发布广播
func Pub(name string, msg interface{}, opts ...client.PublishOption) error {
	return APP().Pub(name, msg, opts...)
}

// Web服务
func Web() *web.Server {
	return APP().Web
}
