package web

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
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

func (c *Context) Session() (*sessions.Session, error) {
	return session.Get("SESSION", c.Ctx())
}

func (c *Context) SessionSave(session *sessions.Session) error {
	return session.Save(c.Ctx().Request(), c.Ctx().Response())
}

func (c *Context) SessionDo(closure func(*sessions.Session) interface{}) (interface{}, error) {
	s, err := c.Session()
	if err != nil {
		return nil, err
	}

	res := closure(s)

	if err := c.SessionSave(s); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Context) SessId() string {
	res, err := c.SessionDo(func(ss *sessions.Session) interface{} {
		return ss.ID
	})
	if err != nil {
		return ""
	}

	return res.(string)
}

func (c *Context) SessOpts() *sessions.Options {
	res, err := c.SessionDo(func(ss *sessions.Session) interface{} {
		return ss.Options
	})
	if err != nil {
		return nil
	}

	return res.(*sessions.Options)
}

func (c *Context) SessOptsSet(opts *sessions.Options) error {
	_, err := c.SessionDo(func(ss *sessions.Session) interface{} {
		ss.Options = opts
		return nil
	})
	return err
}

func (c *Context) SessFlashAdd(val interface{}, key ...string) error {
	_, err := c.SessionDo(func(ss *sessions.Session) interface{} {
		ss.AddFlash(val, key...)
		return nil
	})
	return err
}

func (c *Context) SessFlash(key ...string) []interface{} {
	res, err := c.SessionDo(func(ss *sessions.Session) interface{} {
		return ss.Flashes(key...)
	})
	if err != nil {
		return nil
	}

	return res.([]interface{})
}

func (c *Context) SessAll() map[interface{}]interface{} {
	res, err := c.SessionDo(func(ss *sessions.Session) interface{} {
		return ss.Values
	})
	if err != nil {
		return make(map[interface{}]interface{})
	}

	return res.(map[interface{}]interface{})
}

func (c *Context) SessHas(key interface{}) bool {
	_, b := c.SessAll()[key]
	return b
}

func (c *Context) SessSet(values map[interface{}]interface{}) error {
	_, err := c.SessionDo(func(ss *sessions.Session) interface{} {
		for k, v := range values {
			ss.Values[k] = v
		}
		return nil
	})
	return err
}

func (c *Context) SessGetOne(key interface{}) interface{} {
	return c.SessAll()[key]
}

func (c *Context) SessSetOne(key, val interface{}) error {
	return c.SessSet(map[interface{}]interface{}{
		key: val,
	})
}

func ExtendCtx(ctx echo.Context) *Context {
	c := &Context{
		ctx: ctx,
	}
	return c
}
