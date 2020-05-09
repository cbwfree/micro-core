package srv

import (
	"github.com/micro/cli/v2"
)

type Option func(o *Options)

type Options struct {
	Dev  bool   // 开发模式
	Root string // 数据保存位置

	RedisUrl      string // Redis URL地址
	RedisIdeConns int    // Redis 最小空闲连接数
	RedisMaxPool  int    // Redis 最大连接数

	MongoUrl     string // MongoDB URL地址
	MongoMinPool uint64 // MongoDB 最小连接数
	MongoMaxPool uint64 // MongoDB 最大连接数

	HttpAddr        string           // HTTP 服务地址
	HttpTimeout     int64            // HTTP 请求超时
	HttpStaticUri   string           // HTTP 静态文件服务URI路径
	HttpStaticRoot  string           // HTTP 静态文件服务本地路径
	HttpAllowOrigin *cli.StringSlice // HTTP 允许的跨域源
}

func newOptions() *Options {
	return &Options{
		HttpAllowOrigin: cli.NewStringSlice("*"),
	}
}