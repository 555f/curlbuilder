package curlbuilder

import (
	"net/http"
	"sort"

	"github.com/labstack/echo/v4"
)

func FromRequest(r *http.Request) string {
	keys := make([]string, 0, len(r.Header))
	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	headers := make([]string, 0, len(keys)*2)
	for _, headerName := range keys {
		headers = append(headers, headerName, r.Header.Get(headerName))
	}

	return New().
		SetBody(r.Body).
		SetHeaders(headers...).
		SetMethod(r.Method).
		SetURL(r.URL.String()).
		String()
}

func FromEchoContext(ctx echo.Context) string {
	return FromRequest(ctx.Request())
}
