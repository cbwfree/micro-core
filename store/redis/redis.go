package rds

import (
	"fmt"
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v7"
	log "github.com/micro/go-micro/v2/logger"
	"strings"
)

// Redis 存储
type Store struct {
	db     int
	opts   *Options
	client *redis.Client
	locker *redislock.Client
}

func (rs *Store) With(opts ...Option) {
	rs.opts.With(opts...)
}

func (rs *Store) Opts() *Options {
	return rs.opts
}

func (rs *Store) Locker() *redislock.Client {
	return rs.locker
}

func (rs *Store) Client() *redis.Client {
	return rs.client
}

func (rs *Store) Connect() error {
	if rs.client != nil {
		return nil
	}

	if rs.opts.Uri == "" {
		rs.opts.RawUrl = "redis://127.0.0.1:6379"
	} else if !strings.HasPrefix(rs.opts.Uri, "redis://") {
		rs.opts.RawUrl = "redis://" + rs.opts.Uri
	} else {
		rs.opts.RawUrl = rs.opts.Uri
	}

	opts, err := redis.ParseURL(rs.opts.RawUrl)
	if err != nil {
		return fmt.Errorf("invalid redis url: %s", err.Error())
	}

	// 设置启动参数
	opts.DB = rs.db
	opts.MaxRetries = rs.opts.MaxRetries
	opts.MinRetryBackoff = rs.opts.MinRetryBackoff
	opts.MaxRetryBackoff = rs.opts.MaxRetryBackoff
	opts.ReadTimeout = rs.opts.ReadTimeout
	opts.WriteTimeout = rs.opts.WriteTimeout
	opts.PoolSize = rs.opts.PoolSize
	opts.MinIdleConns = rs.opts.MinIdleConns
	opts.IdleTimeout = rs.opts.IdleTimeout

	rs.client = redis.NewClient(opts)
	if err := rs.client.Ping().Err(); err != nil {
		return err
	}

	// 启用分布式锁
	rs.locker = redislock.New(rs.client)

	log.Debugf("Store [redis] Connect to %s", opts.Addr)

	return nil
}

func (rs *Store) Disconnect() error {
	if rs.client != nil {
		if err := rs.client.Close(); err != nil {
			return err
		}
		rs.client = nil
	}
	return nil
}

// Do 执行命令
func (rs *Store) Do(args ...interface{}) *redis.Cmd {
	cmd := redis.NewCmd(args...)
	_ = rs.client.Process(cmd)
	return cmd
}

func (rs *Store) Pipelined(fn func(pipe redis.Pipeliner) error) ([]redis.Cmder, error) {
	return rs.client.Pipelined(fn)
}

func (rs *Store) TxPipelined(fn func(tx redis.Pipeliner) error) ([]redis.Cmder, error) {
	return rs.client.TxPipelined(fn)
}

// HGetStruct HGETALL 转结构体
func (rs *Store) HGetStruct(key string, result interface{}) error {
	cmd := redis.NewSliceCmd(CmdHGetAll, key)
	_ = rs.client.Process(cmd)
	res, err := cmd.Result()
	if err != nil {
		return err
	}
	return ScanStruct(res, result)
}

// HSetStruct 设置结构体
func (rs *Store) HSetStruct(key string, result interface{}) error {
	cmd := redis.NewStatusCmd(Args{CmdHSet}.Add(key).AddFlat(result)...)
	_ = rs.client.Process(cmd)
	_, err := cmd.Result()
	return err
}

func NewStore(db int, opts ...Option) *Store {
	rs := &Store{
		db:   db,
		opts: newOptions(opts...),
	}
	return rs
}
