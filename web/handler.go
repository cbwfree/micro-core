package web

import (
	"github.com/labstack/echo/v4"
)

// 统一错误处理
func errorHandler(err error, ctx echo.Context) {
	res := ParseError(err)

	// Send response
	if !ctx.Response().Committed {
		if ctx.Request().Method == echo.HEAD { // Issue #608
			err = ctx.NoContent(res.Code)
		} else {
			err = CtxResult(ctx, res)
		}
		if err != nil {
			ctx.Logger().Error(err)
		}
	}
}
