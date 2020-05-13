package web

import (
	"github.com/cbwfree/micro-core/conv"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/steambap/captcha"
	"image/color"
)

type Context struct {
	ctx echo.Context
}

func (c *Context) Ctx() echo.Context {
	return c.ctx
}

func (c *Context) Error(err error) error {
	return CtxResult(c.ctx, ParseError(err))
}

func (c *Context) JsonError(code int, msg ...string) error {
	return CtxError(c.ctx, code, msg...)
}

func (c *Context) JsonSuccess(data interface{}, msg ...string) error {
	return CtxSuccess(c.ctx, data, msg...)
}

func (c *Context) Bind(req interface{}) error {
	return c.ctx.Bind(req)
}

func (c *Context) BindValid(req interface{}) error {
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.ctx.Validate(req); err != nil {
		return err
	}
	return nil
}

func (c *Context) Session() *sessions.Session {
	s, _ := session.Get("SESSION", c.Ctx())
	return s
}

func (c *Context) SessionDo(closure func(*sessions.Session) interface{}) interface{} {
	s := c.Session()
	res := closure(s)
	_ = s.Save(c.ctx.Request(), c.ctx.Response())
	return res
}

func (c *Context) SessId() string {
	res := c.SessionDo(func(ss *sessions.Session) interface{} {
		return ss.ID
	})
	return res.(string)
}

func (c *Context) SessOpts() *sessions.Options {
	res := c.SessionDo(func(ss *sessions.Session) interface{} {
		return ss.Options
	})
	return res.(*sessions.Options)
}

func (c *Context) SessOptsSet(opts *sessions.Options) {
	_ = c.SessionDo(func(ss *sessions.Session) interface{} {
		ss.Options = opts
		return nil
	})
}

func (c *Context) SessFlashAdd(val interface{}, key ...string) {
	_ = c.SessionDo(func(ss *sessions.Session) interface{} {
		ss.AddFlash(val, key...)
		return nil
	})
}

func (c *Context) SessFlash(key ...string) []interface{} {
	res := c.SessionDo(func(ss *sessions.Session) interface{} {
		return ss.Flashes(key...)
	})
	return res.([]interface{})
}

func (c *Context) SessAll() map[interface{}]interface{} {
	res := c.SessionDo(func(ss *sessions.Session) interface{} {
		return ss.Values
	})
	return res.(map[interface{}]interface{})
}

func (c *Context) SessHas(key interface{}) bool {
	_, b := c.SessAll()[key]
	return b
}

func (c *Context) SessSet(values map[interface{}]interface{}) {
	_ = c.SessionDo(func(ss *sessions.Session) interface{} {
		for k, v := range values {
			ss.Values[k] = v
		}
		return nil
	})
}

func (c *Context) SessGetOne(key interface{}) interface{} {
	return c.SessAll()[key]
}

func (c *Context) SessSetOne(key, val interface{}) {
	c.SessSet(map[interface{}]interface{}{
		key: val,
	})
}

// 创建验证码
func (c *Context) CaptchaNew(key string, width, height int, setOpt ...captcha.SetOption) error {
	var opt captcha.SetOption
	if len(setOpt) > 0 && setOpt[0] != nil {
		opt = setOpt[0]
	} else {
		opt = func(opt *captcha.Options) {
			opt.BackgroundColor = color.White
			opt.CharPreset = "0123456789"
			opt.CurveNumber = 1
			opt.FontDPI = 80
		}
	}

	data, err := captcha.New(width, height, opt)
	if err != nil {
		return c.Error(err)
	}

	c.SessSetOne(key, data.Text)

	return data.WriteImage(c.ctx.Response().Writer)
}

// 验证验证码
func (c *Context) CaptchaCheck(key string, captcha string) bool {
	codeText := c.SessionDo(func(s *sessions.Session) interface{} {
		if c, ok := s.Values[key]; ok {
			return c
		}
		return ""
	})
	return conv.String(codeText) == captcha
}

func ExtendCtx(ctx echo.Context) *Context {
	c := &Context{
		ctx: ctx,
	}
	return c
}
