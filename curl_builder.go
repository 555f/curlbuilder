package curlbuilder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type CurlBuilder struct {
	url     string
	https   bool
	body    interface{}
	method  string
	headers []string
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
		headers[key] = append(headers[key], value)
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		_, _ = fmt.Fprintf(buf, "-H %s: %s ", key, escape(strings.Join(headers[key], " ")))
	}

	_, _ = fmt.Fprintf(buf, b.url)

	return buf.String()
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

func New() *CurlBuilder {
	return &CurlBuilder{
		method: http.MethodGet,
	}
}
