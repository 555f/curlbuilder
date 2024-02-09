package curlbuilder

import (
	"bytes"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

func FromRequest(r *http.Request) string {
	return New().SetRequest(r).String()
}

func FromEchoContext(ctx echo.Context) string {
	return FromRequest(ctx.Request())
}

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == nil || b == http.NoBody {
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
