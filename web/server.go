package web

import (
	"encoding/json"
	"github.com/cbwfree/micro-core/conv"
	"github.com/cbwfree/micro-core/fn"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/micro/go-micro/v2/logger"
	"io/ioutil"
	"net"
	"path/filepath"
	"strings"
	"sync"
)

type Server struct {
	sync.Mutex
	running bool
	exit    chan chan error

	echo *echo.Echo
	opts *Options

	socket *Socket
}

func (s *Server) Echo() *echo.Echo {
	return s.echo
}

func (s *Server) Socket() *Socket {
	return s.socket
}

func (s *Server) With(opts ...Option) {
	s.opts.With(opts...)
}

func (s *Server) Opts() *Options {
	return s.opts
}

// 启用API路由
func (s *Server) enableAPIRoutes() {
	if s.opts.APIPrefix == "" {
		return
	}

	api := s.echo.Group(s.opts.APIPrefix)

	for _, r := range s.opts.APIRoutes {
		r(api)
	}
}

// 启用跨域
func (s *Server) enableCORS() {
	if s.opts.AllowOrigins == "" {
		return
	}

	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{s.opts.AllowOrigins},
		AllowCredentials: true,
		AllowHeaders:     s.opts.AllowHeaders,
		AllowMethods:     s.opts.AllowMethods,
		ExposeHeaders:    s.opts.ExposeHeaders,
	}))
}

// 启用静态文件
func (s *Server) enableStatic() {
	if len(s.opts.StaticRoot) == 0 {
		return
	}

	for _, val := range strings.Split(s.opts.StaticRoot, ":") {
		prefix := "/"
		root := val
		confFile := filepath.Join(val, "config.json")

		// 读取配置文件
		if fn.FileExist(confFile) {
			b, err := ioutil.ReadFile(confFile)
			if err != nil {
				log.Fatal("Read web static config.json error: %v", err)
				return
			}

			var res = make(map[string]interface{})
			if err := json.Unmarshal(b, &res); err != nil {
				log.Fatal("Parse web static config.json error: %v", err)
				return
			}

			prefix = conv.String(res["baseUrl"])
		}

		var use func(mdd ...echo.MiddlewareFunc)
		if prefix != "/" {
			use = s.echo.Group(prefix).Use
		} else {
			use = s.echo.Use
		}

		use(middleware.StaticWithConfig(middleware.StaticConfig{
			Root:   root,
			Index:  "index.html",
			HTML5:  true,  // SPA 单页面是否转发
			Browse: false, // 是否启用目录浏览
		}))

		log.Infof("HTTP Server Enable Static Service, Prefix: %s, Path: %s", prefix, root)
	}
}

// 启用Session
func (s *Server) enableSession() {
	if s.opts.SessionStore == "" {
		return
	}

	// 检查目录是否存在
	store := filepath.Join(s.opts.Root, s.opts.SessionStore)
	if err := fn.Mkdir(store); err != nil {
		log.Fatalf("Enable Web Session Error: %s", err)
		return
	}

	s.echo.Use(session.Middleware(
		sessions.NewFilesystemStore(store, []byte(s.opts.SessionSecret)),
	))

	log.Infof("HTTP Server Enable Session Service, Save Path: %s", store)
}

// 启用WebSocket
func (s *Server) enableSocket() {
	if s.opts.SocketPath == "" {
		return
	}

	s.socket = NewSocket(s)
	s.echo.GET(s.opts.SocketPath, s.socket.Handler)
}

func (s *Server) Start() error {
	s.Lock()
	defer s.Unlock()

	if s.running {
		return nil
	}

	s.enableCORS()      // 启用跨域
	s.enableSession()   // 启用Session
	s.enableSocket()    // 启用WebSocket
	s.enableAPIRoutes() // 注册API路由
	s.enableStatic()    // 启用静态文件

	l, err := net.Listen("tcp", s.opts.Addr)
	if err != nil {
		return err
	}

	s.echo.Listener = l

	go func() {
		_ = s.echo.Start("")
	}()

	s.exit = make(chan chan error, 1)
	s.running = true

	go func() {
		ch := <-s.exit
		ch <- l.Close()
	}()

	log.Infof("HTTP Server Listening on %v", l.Addr().String())

	return nil
}

// 关闭服务
func (s *Server) Close() error {
	s.Lock()
	defer s.Unlock()

	if !s.running {
		return nil
	}

	ch := make(chan error, 1)
	s.exit <- ch
	s.running = false

	log.Infof("HTTP Server Close ... ")

	return <-ch
}

func NewServer(opts ...Option) *Server {
	s := &Server{
		echo: echo.New(),
		opts: newOptions(opts...),
	}

	s.echo.HideBanner = true // 隐藏Echo的Banner
	s.echo.HidePort = true
	s.echo.HTTPErrorHandler = errorHandler // 统一错误处理
	s.echo.Validator = NewWebValidator()   // 数据验证器

	// 记录请求日志
	s.echo.Use(middleware.Logger())
	// 从 panic 链中的任意位置恢复程序， 打印堆栈的错误信息，并将错误集中交给 HTTPErrorHandler 处理。
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.Secure())

	return s
}
