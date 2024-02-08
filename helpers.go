package curlbuilder

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func FromRequest(r *http.Request) string {
	return New().SetRequest(r).String()
}

func FromEchoContext(ctx echo.Context) string {
	return FromRequest(ctx.Request())
}
