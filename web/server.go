package web

import (
	"github.com/cbwfree/micro-core/fn"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/micro/go-micro/v2/util/log"
	"net"
	"sync"
)

type Server struct {
	sync.Mutex
	name    string
	running bool
	exit    chan chan error

	echo *echo.Echo
	opts *Options
}

func (s *Server) Echo() *echo.Echo {
	return s.echo
}

func (s *Server) With(opts ...Option) {
	s.opts.With(opts...)
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
	if len(s.opts.AllowOrigins) == 0 {
		return
	}

	s.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  s.opts.AllowOrigins,
		AllowHeaders:  s.opts.AllowHeaders,
		AllowMethods:  s.opts.AllowMethods,
		ExposeHeaders: s.opts.ExposeHeaders,
	}))
}

// 启用静态文件
func (s *Server) enableStatic() {
	if s.opts.StaticRoot == "" {
		return
	}

	var use func(mdd ...echo.MiddlewareFunc)
	if s.opts.StaticUri == "" {
		static := s.echo.Group(s.opts.StaticUri)
		use = static.Use
	} else {
		use = s.echo.Use
	}

	use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   s.opts.StaticRoot,
		Index:  "index.html",
		HTML5:  true,  // SPA 单页面是否转发
		Browse: false, // 是否启用目录浏览
	}))

	log.Info("[%s] HTTP Server Enable Static Service, Link: %s", s.name, s.opts.StaticUri)
}

// 启用Session
func (s *Server) enableSession() {
	if !s.opts.EnableSession {
		return
	}

	// 检查目录是否存在
	if !fn.ExistDir(s.opts.SessionStore) {
		if err := fn.MkDir(s.opts.SessionStore, 0755); err != nil {
			log.Fatal("Enable Web Session Error: %s", err)
			return
		}
	}

	s.echo.Use(session.Middleware(
		sessions.NewFilesystemStore(
			s.opts.SessionStore,
			[]byte(DefaultSessionSecret),
		),
	))

	log.Info("[%s] HTTP Server Enable Session Service, Save: %s", s.name, s.opts.SessionStore)
}

func (s *Server) Start() error {
	s.Lock()
	defer s.Unlock()

	if s.running {
		return nil
	}

	s.enableCORS()      // 启用跨域
	s.enableSession()   // 启用Session
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

	log.Info("[%s] HTTP Server Listening on %v", s.name, l.Addr().String())

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

	log.Info("[%s] HTTP Server Close ... ", s.name)

	return <-ch
}

func NewServer(name string, opts ...Option) *Server {
	s := &Server{
		name: name,
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
	s.echo.Use(middleware.CSRF())

	return s
}
