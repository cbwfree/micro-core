package srv

import (
	"context"
	"fmt"
	"github.com/cbwfree/micro-core/compile"
	"github.com/cbwfree/micro-core/fn"
	mgo "github.com/cbwfree/micro-core/store/mongo"
	rds "github.com/cbwfree/micro-core/store/redis"
	"github.com/cbwfree/micro-core/web"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/selector"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server"
	"strings"
	"sync"
	"time"
)

var (
	app = &App{
		opts:      newOptions(),
		publisher: make(map[string]micro.Publisher),
	}
)

type App struct {
	sync.Mutex
	srv    micro.Service
	ctx    context.Context
	cancel context.CancelFunc
	opts   *Options

	Mongo *mgo.Store
	Redis *rds.Store
	Web   *web.Server

	publisher map[string]micro.Publisher // 订阅
}

func (a *App) With(with ...WithAPP) {
	for _, w := range with {
		w(a)
	}
}

func (a *App) New(srvName string, flags ...[]cli.Flag) {
	if a.srv != nil {
		return
	}

	compile.SetName(srvName)

	// 创建服务
	a.ctx, a.cancel = context.WithCancel(context.Background())
	a.srv = micro.NewService(
		micro.Name(compile.Name()),
		micro.Version(compile.Version()),
		micro.Context(a.ctx),
		micro.Flags(mergeFlags(flags)...),
		micro.BeforeStart(func() error {
			// 启动时输出版本信息
			compile.EchoVersion(a.srv)

			if a.opts.Dev {
				log.Warnf("the current startup mode is development mode")
			}

			// 检查目录是否存在
			if err := fn.Mkdir(a.opts.Root, 0755); err != nil {
				return err
			}

			if a.Redis != nil {
				a.Redis.Opts().Uri = a.opts.RedisUrl
				a.Redis.Opts().Db = a.opts.RedisDb
				a.Redis.Opts().MinIdleConns = a.opts.RedisIdeConns
				a.Redis.Opts().PoolSize = a.opts.RedisMaxPool
				if err := a.Redis.Connect(); err != nil {
					return err
				}
			}

			if a.Mongo != nil {
				if a.opts.MongoDb == "" {
					a.Mongo.Opts().Db = strings.Replace(srvName, ".", "-", -1)
				} else {
					a.Mongo.Opts().Db = strings.Replace(a.opts.MongoDb, ".", "-", -1)
				}
				a.Mongo.Opts().Uri = a.opts.MongoUrl
				a.Mongo.Opts().MinPoolSize = a.opts.MongoMinPool
				a.Mongo.Opts().MaxPoolSize = a.opts.MongoMaxPool
				if err := a.Mongo.Connect(); err != nil {
					return err
				}
			}

			return nil
		}),
		micro.AfterStart(func() error {
			if a.Web != nil {
				a.Web.Opts().Addr = a.opts.HttpAddr
				a.Web.Opts().Timeout = time.Duration(a.opts.HttpTimeout) * time.Second
				a.Web.Opts().Root = a.opts.Root
				a.Web.Opts().StaticRoot = a.opts.HttpStaticRoot
				a.Web.Opts().AllowOrigins = a.opts.HttpAllowOrigin
				if err := a.Web.Start(); err != nil {
					return err
				}
			}
			return nil
		}),
		micro.BeforeStop(func() error {
			if a.Web != nil {
				if err := a.Web.Close(); err != nil {
					return err
				}
			}

			return nil
		}),
		micro.AfterStop(func() error {
			if a.Redis != nil {
				if err := a.Redis.Disconnect(); err != nil {
					return err
				}
			}

			if a.Mongo != nil {
				if err := a.Mongo.Disconnect(); err != nil {
					return err
				}
			}

			return nil
		}),
	)
}

// Close 主动关闭APP
func (a *App) Close() {
	a.cancel()
}

func (a *App) Init(opts ...micro.Option) {
	a.srv.Init(opts...)
}

// Service 获取服务对象
func (a *App) Srv() micro.Service {
	return a.srv
}

// Version 获取服务版本
func (a *App) Version() string {
	return a.srv.Server().Options().Version
}

// Id 获取服务ID
func (a *App) Id() string {
	return a.srv.Server().Options().Id
}

// Name 获取服务名称
func (a *App) Name() string {
	return a.srv.Server().Options().Name
}

// NameId 获取服务节点ID
func (a *App) NameId() string {
	return fmt.Sprintf("%s-%s", a.Name(), a.Id())
}

// Server 获取服务的服务端
func (a *App) SrvServer() server.Server {
	return a.srv.Server()
}

// Client 获取服务的客户端
func (a *App) SrvClient() client.Client {
	return a.srv.Client()
}

// 获取服务节点列表
func (a *App) GetServices(name string) []*registry.Service {
	res, _ := a.srv.Options().Registry.GetService(name)
	return res
}

// 获取服务节点列表
func (a *App) ListServices() []*registry.Service {
	res, _ := a.srv.Options().Registry.ListServices()
	return res
}

// PublishCtx 发布消息 (自定义Context信息)
func (a *App) PubCtx(ctx context.Context, name string, msg interface{}, opts ...client.PublishOption) error {
	if pub, ok := a.publisher[name]; ok {
		return pub.Publish(ctx, msg, opts...)
	}
	return fmt.Errorf("not found [%s] publisher", name)
}

// Publish 发布消息
func (a *App) Pub(name string, msg interface{}, opts ...client.PublishOption) error {
	return a.PubCtx(context.TODO(), name, msg, opts...)
}

// Call 通过名称调用RPC
// 	@name 服务名称
// 	@method rpc方法名称. 即 serviceName.rpcName
// 	@in 请求参数
// 	@out 返回数据
// 	@filter 节点选择
func (a *App) Call(name string, method string, in interface{}, out interface{}, filter ...selector.Filter) error {
	return a.CallCtx(context.TODO(), name, method, in, out, filter...)
}

// CallCtx 通过名称调用RPC
// 	@name 服务名称
// 	@method rpc方法名称. 即 serviceName.rpcName
// 	@in 请求参数
// 	@out 返回数据
// 	@filter 节点选择
func (a *App) CallCtx(ctx context.Context, name string, method string, in interface{}, out interface{}, filter ...selector.Filter) error {
	var opts []client.CallOption
	if len(filter) > 0 {
		opts = append(opts, FilterSelector(filter[0]))
	}
	req := a.srv.Client().NewRequest(name, method, in)
	return a.srv.Client().Call(ctx, req, out, opts...)
}
