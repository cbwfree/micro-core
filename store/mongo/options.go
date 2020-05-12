package mgo

import "time"

const (
	DefaultMinPoolSize      uint64 = 20              // 最小连接池大小
	DefaultMaxPoolSize      uint64 = 100             // 最大连接池大小
	DefaultConnectTimeout          = 5 * time.Second // 连接超时时间
	DefaultSocketTimeout           = 5 * time.Second
	DefaultMaxConnIdleTime         = 3 * time.Second // 最大空闲时间
	DefaultReadWriteTimeout        = 3 * time.Second // 读写超时时间
)

type Option func(o *Options)

type Options struct {
	Uri              string
	RawUrl           string
	Db               string
	MinPoolSize      uint64
	MaxPoolSize      uint64
	ConnectTimeout   time.Duration
	SocketTimeout    time.Duration
	MaxConnIdleTime  time.Duration
	ReadWriteTimeout time.Duration
}

func (o *Options) With(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func newOptions(opts ...Option) *Options {
	o := &Options{
		MinPoolSize:      DefaultMinPoolSize,
		MaxPoolSize:      DefaultMaxPoolSize,
		ConnectTimeout:   DefaultConnectTimeout,
		SocketTimeout:    DefaultSocketTimeout,
		MaxConnIdleTime:  DefaultMaxConnIdleTime,
		ReadWriteTimeout: DefaultReadWriteTimeout,
	}
	o.With(opts...)
	return o
}

func WithConnectTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.ConnectTimeout = t
	}
}

func WithSocketTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.SocketTimeout = t
	}
}

func WithMaxConnIdleTime(t time.Duration) Option {
	return func(o *Options) {
		o.MaxConnIdleTime = t
	}
}

func WithReadWriteTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.ReadWriteTimeout = t
	}
}
