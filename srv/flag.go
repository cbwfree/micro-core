package srv

import (
	"github.com/micro/cli/v2"
	"os"
)

var (
	FlagBasic = []cli.Flag{
		&cli.BoolFlag{
			Name:        "dev",
			Usage:       "设置当前运行模式为开发模式",
			Destination: &OPTS().Dev,
		},
		&cli.StringFlag{
			Name:        "root",
			Value:       os.TempDir(),
			Usage:       "设置数据保存位置",
			EnvVars:     []string{"CORE_ROOT"},
			Destination: &OPTS().Root,
		},
	}

	FlagRedis = []cli.Flag{
		&cli.StringFlag{
			Name:        "redis",
			Value:       "",
			Usage:       "设置redis连接地址. 格式: redis://[:password@]hostname:port/[db]",
			EnvVars:     []string{"CORE_REDIS_URL"},
			Destination: &OPTS().RedisUrl,
		},
		&cli.IntFlag{
			Name:        "redis_db",
			Value:       0,
			Usage:       "设置redis连接的db",
			EnvVars:     []string{"CORE_REDIS_DB"},
			Destination: &OPTS().RedisDb,
		},
		&cli.IntFlag{
			Name:        "redis_ide_conns",
			Value:       50,
			Usage:       "设置redis最大空闲连接数",
			EnvVars:     []string{"CORE_REDIS_IDE_CONNS"},
			Destination: &OPTS().RedisIdeConns,
		},
		&cli.IntFlag{
			Name:        "redis_max_pool",
			Value:       100,
			Usage:       "设置redis最大连接数",
			EnvVars:     []string{"CORE_REDIS_MAX_POOL"},
			Destination: &OPTS().RedisMaxPool,
		},
	}

	FlagMongo = []cli.Flag{
		&cli.StringFlag{
			Name:        "mongo",
			Value:       "",
			Usage:       "设置MongoDB连接地址. 格式: mongodb://[username:password@]host1[:port1][,host2[:port2],...[,hostN[:portN]]]",
			EnvVars:     []string{"CORE_MONGO_URL"},
			Destination: &OPTS().MongoUrl,
		},
		&cli.StringFlag{
			Name:        "mongo_db",
			Value:       "",
			Usage:       "设置MongoDB连接的默认数据库",
			EnvVars:     []string{"CORE_MONGO_DB"},
			Destination: &OPTS().MongoDb,
		},
		&cli.Uint64Flag{
			Name:        "mongo_min_pool",
			Value:       40,
			Usage:       "设置MongoDB最小连接数",
			EnvVars:     []string{"CORE_MONGO_MIN_POOL"},
			Destination: &OPTS().MongoMinPool,
		},
		&cli.Uint64Flag{
			Name:        "mongo_max_pool",
			Value:       100,
			Usage:       "设置MongoDB最大连接数",
			EnvVars:     []string{"CORE_MONGO_MAX_POOL"},
			Destination: &OPTS().MongoMaxPool,
		},
	}

	FlagHttp = []cli.Flag{
		&cli.StringFlag{
			Name:        "http_addr",
			Value:       "",
			Usage:       "设置HTTP服务监听地址.",
			EnvVars:     []string{"CORE_HTTP_ADDR"},
			Destination: &OPTS().HttpAddr,
		},
		&cli.Int64Flag{
			Name:        "http_timeout",
			Value:       30,
			Usage:       "设置HTTP服务请求超时时间",
			EnvVars:     []string{"CORE_HTTP_TIMEOUT"},
			Destination: &OPTS().HttpTimeout,
		},
		&cli.StringFlag{
			Name:        "http_static_uri",
			Value:       "",
			Usage:       "设置HTTP静态文件目录访问URI. 多个目录使用:分隔",
			EnvVars:     []string{"CORE_HTTP_STATIC_URI"},
			Destination: &OPTS().HttpStaticUri,
		},
		&cli.StringFlag{
			Name:        "http_static_root",
			Value:       "",
			Usage:       "设置HTTP静态文件目录路径. 多个目录使用:分隔",
			EnvVars:     []string{"CORE_HTTP_STATIC_ROOT"},
			Destination: &OPTS().HttpStaticRoot,
		},
		&cli.StringFlag{
			Name:        "http_allow_origins",
			Value:       "*",
			Usage:       "设置HTTP运行的跨域源地址.",
			EnvVars:     []string{"CORE_HTTP_ALLOW_ORIGINS"},
			Destination: &OPTS().HttpAllowOrigin,
		},
	}
)

func mergeFlags(flags [][]cli.Flag) []cli.Flag {
	var result []cli.Flag
	for _, flag := range flags {
		result = append(result, flag...)
	}
	return result
}
