package rds

import "time"

const (
	DefaultMinIdleConns    = 20              // 最小空闲连接数
	DefaultPoolSize        = 100             // 最大连接数
	DefaultMaxRetries      = 3               // 最大重试次数
	DefaultReadTimeout     = 3 * time.Second // 读取超时
	DefaultWriteTimeout    = 3 * time.Second // 写入超时
	DefaultIdleTimeout     = 5 * time.Minute // 空闲连接关闭时间
	DefaultMinRetryBackoff = 200 * time.Millisecond
	DefaultMaxRetryBackoff = 1 * time.Second
)

type Option func(o *Options)

type Options struct {
	Uri             string
	RawUrl          string
	Db              int
	MinIdleConns    int
	PoolSize        int
	MaxRetries      int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration
}

func (o *Options) With(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func newOptions(opts ...Option) *Options {
	o := &Options{
		MinIdleConns:    DefaultMinIdleConns,
		PoolSize:        DefaultPoolSize,
		MaxRetries:      DefaultMaxRetries,
		ReadTimeout:     DefaultReadTimeout,
		WriteTimeout:    DefaultWriteTimeout,
		IdleTimeout:     DefaultIdleTimeout,
		MinRetryBackoff: DefaultMinRetryBackoff,
		MaxRetryBackoff: DefaultMaxRetryBackoff,
	}
	o.With(opts...)
	return o
}

func WithMaxRetries(size int) Option {
	return func(o *Options) {
		o.MaxRetries = size
	}
}

func WithReadTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.ReadTimeout = t
	}
}

func WithWriteTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.WriteTimeout = t
	}
}

func WithIdleTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.IdleTimeout = t
	}
}

func WithMinRetryBackoff(t time.Duration) Option {
	return func(o *Options) {
		o.MinRetryBackoff = t
	}
}

func WithMaxRetryBackoff(t time.Duration) Option {
	return func(o *Options) {
		o.MaxRetryBackoff = t
	}
}
