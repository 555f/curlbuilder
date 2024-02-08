package curlbuilder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
)

type CurlBuilder struct {
	url     string
	body    interface{}
	method  string
	headers []string
	secrets map[string]struct{}
}

func (b *CurlBuilder) String() string {
	buf := bytes.NewBuffer(nil)
	_, _ = fmt.Fprintf(buf, "curl ")

	if strings.HasPrefix(b.url, "https://") {
		_, _ = fmt.Fprintf(buf, "-k ")
	}

	_, _ = fmt.Fprintf(buf, "-X %s ", b.method)

	if b.body != nil {
		var body []byte
		switch t := b.body.(type) {
		default:
			body, _ = json.Marshal(b.body)
		case io.Reader:
			body, _ = io.ReadAll(t)
		case string:
			body = []byte(t)
		case []byte:
			body = t
		}
		if len(body) > 0 {
			_, _ = fmt.Fprintf(buf, "-d %s ", escape(string(body)))
		}
	}

	var (
		headers = map[string][]string{}
		keys    = make([]string, 0, len(b.headers))
	)

	for i := 0; i < len(b.headers); i += 2 {
		key := b.headers[i]
		value := b.headers[i+1]
		if _, ok := b.secrets[key]; ok {
			value = strings.Repeat("*", len(value))
		}
		headers[key] = append(headers[key], value)
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		_, _ = fmt.Fprintf(buf, "-H %s ", escape(key+": "+strings.Join(headers[key], " ")))
	}

	_, _ = fmt.Fprintf(buf, b.url)

	return buf.String()
}

func (b *CurlBuilder) SetRequest(r *http.Request) *CurlBuilder {
	keys := make([]string, 0, len(r.Header))
	for k := range r.Header {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	headers := make([]string, 0, len(keys)*2)
	for _, headerName := range keys {
		headers = append(headers, headerName, r.Header.Get(headerName))
	}
	b.SetBody(r.Body).
		SetHeaders(headers...).
		SetMethod(r.Method).
		SetURL(r.URL.String())
	return b
}

func (b *CurlBuilder) SetEchoContext(ctx echo.Context) *CurlBuilder {
	return b.SetRequest(ctx.Request())
}

func (b *CurlBuilder) SetURL(url string) *CurlBuilder {
	b.url = url
	return b
}

func (b *CurlBuilder) SetHeaders(headers ...string) *CurlBuilder {
	if len(headers)%2 != 0 {
		panic("SetHeaders: headers must be key/value")
	}
	b.headers = headers
	return b
}

func (b *CurlBuilder) SetMethod(method string) *CurlBuilder {
	b.method = method
	return b
}

func (b *CurlBuilder) SetBody(body interface{}) *CurlBuilder {
	b.body = body
	return b
}

func (b *CurlBuilder) SetSecret(fields ...string) *CurlBuilder {
	for _, f := range fields {
		b.secrets[f] = struct{}{}
	}
	return b
}

func New() *CurlBuilder {
	return &CurlBuilder{
		method:  http.MethodGet,
		secrets: make(map[string]struct{}, 32),
	}
}
