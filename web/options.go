package web

import (
	"github.com/labstack/echo/v4"
	"path/filepath"
	"time"
)

const (
	DefaultSessionSecret = "SessionSecret" // 修改Session数据保存秘钥
)

var (
	defaultAllowOrigins  = []string{"*"}
	defaultAllowMethods  = []string{echo.GET, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS}
	defaultAllowHeaders  = []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderXCSRFToken}
	defaultExposeHeaders = []string{echo.HeaderAuthorization, echo.HeaderVary, echo.HeaderCookie}
)

type Route func(e *echo.Group)

type Option func(o *Options)

type Options struct {
	Addr          string
	Timeout       time.Duration
	EnableSocket  bool
	EnableSession bool
	SessionStore  string

	StaticUri     string
	StaticRoot    string
	AllowOrigins  []string
	AllowMethods  []string
	AllowHeaders  []string
	ExposeHeaders []string

	APIPrefix string
	APIRoutes []Route
}

func (o *Options) With(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func newOptions(opts ...Option) *Options {
	o := &Options{
		AllowOrigins:  defaultAllowOrigins,
		AllowMethods:  defaultAllowMethods,
		AllowHeaders:  defaultAllowHeaders,
		ExposeHeaders: defaultExposeHeaders,
	}
	o.With(opts...)
	return o
}

func WithAddr(addr string) Option {
	return func(o *Options) {
		o.Addr = addr
	}
}

func WithTimeout(timeout int64) Option {
	return func(o *Options) {
		o.Timeout = time.Duration(timeout) * time.Second
	}
}

func WithEnableSocket(enable bool) Option {
	return func(o *Options) {
		o.EnableSocket = enable
	}
}

func WithEnableSession(enable bool) Option {
	return func(o *Options) {
		o.EnableSession = enable
	}
}

func WithSessionStore(root string) Option {
	return func(o *Options) {
		if root != "" {
			o.SessionStore = filepath.Join(root, "_session")
		}
	}
}

func WithStaticUri(uri string) Option {
	return func(o *Options) {
		o.StaticUri = uri
	}
}

func WithStaticRoot(static string) Option {
	return func(o *Options) {
		o.StaticRoot = static
	}
}

func WithAllowOrigins(allow []string) Option {
	return func(o *Options) {
		o.AllowOrigins = allow
	}
}

func WithAllowMethods(allow []string) Option {
	return func(o *Options) {
		o.AllowMethods = allow
	}
}

func WithAllowHeaders(allow []string) Option {
	return func(o *Options) {
		o.AllowHeaders = allow
	}
}

func WithExposeHeaders(expose []string) Option {
	return func(o *Options) {
		o.ExposeHeaders = expose
	}
}

func WithAPIPrefix(prefix string) Option {
	return func(o *Options) {
		o.APIPrefix = prefix
	}
}

func WithAPIRoutes(routes ...Route) Option {
	return func(o *Options) {
		o.APIRoutes = append(o.APIRoutes, routes...)
	}
}
