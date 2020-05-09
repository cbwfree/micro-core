# micro-core

## Example

```go
package main

import (
	"github.com/cbwfree/micro-core/srv"
	"github.com/cbwfree/micro-core/web"
	log "github.com/micro/go-micro/v2/logger"
)

func main() {
	srv.New(
		"admin",
		srv.FlagBasic,
		srv.FlagRedis,
		srv.FlagMongo,
		srv.FlagHttp,
	)

	srv.With(
		srv.WithMongoDB("admin"),
		srv.WithRedisDB(0),
		srv.WithWebServer("admin"),
	)

	srv.Web().With(
		web.WithEnableSession(true),
		web.WithAPIPrefix("/api"),
		web.WithAPIRoutes(
            ...
		),
	)

	// 服务初始化
	srv.Init()

	// 启动服务
	if err := srv.Run(); err != nil {
		log.Fatalf("Run Error: %v", err)
	}
}

```