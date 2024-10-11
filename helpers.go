package curlbuilder

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

var _ http.RoundTripper = &RoundTripper{}

type Printer interface {
	Print(ctx context.Context, s string)
}

type RoundTripper struct {
	printer Printer
	next    http.RoundTripper
	secrets []string
}

// RoundTrip implements http.RoundTripper.
func (r *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r.printer.Print(req.Context(), FromRequest(req).SetSecret(r.secrets...).String())
	return r.next.RoundTrip(req)
}

func FromRoundTripper(next http.RoundTripper, printer Printer, secrets ...string) *RoundTripper {
	return &RoundTripper{next: next, printer: printer, secrets: secrets}
}

func FromRequest(r *http.Request) *CurlBuilder {
	return New().SetRequest(r)
}

func FromEchoContext(ctx echo.Context) *CurlBuilder {
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
