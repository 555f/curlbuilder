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

type FormType string

const (
	MultipartFormType  FormType = "multipart"
	UrlencodedFormType FormType = "urlencoded"
)

type CurlBuilder struct {
	url                 string
	body                interface{}
	method              string
	formType            FormType
	headers, formValues []string
	secrets             map[string]struct{}
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
		headers    = map[string][]string{}
		headerKeys = make([]string, 0, len(b.headers))

		formValues = map[string]string{}
		formKeys   = make([]string, 0, len(b.formValues))
	)

	for i := 0; i < len(b.headers); i += 2 {
		key := b.headers[i]
		value := b.headers[i+1]
		if _, ok := b.secrets[key]; ok {
			value = "*****"
		}
		headers[key] = append(headers[key], value)
		headerKeys = append(headerKeys, key)
	}
	sort.Strings(headerKeys)

	for i := 0; i < len(b.formValues); i += 2 {
		key := b.formValues[i]
		value := b.formValues[i+1]
		if _, ok := b.secrets[key]; ok {
			value = strings.Repeat("*", len(value))
		}
		formValues[key] = value
		formKeys = append(formKeys, key)
	}
	sort.Strings(formKeys)

	for _, key := range headerKeys {
		_, _ = fmt.Fprintf(buf, "-H %s ", escape(key+": "+strings.Join(headers[key], " ")))
	}

	if len(formKeys) > 0 {
		switch b.formType {
		case UrlencodedFormType:
			_, _ = fmt.Fprintf(buf, "-d \"")
			for i, key := range formKeys {
				if i > 0 {
					_, _ = fmt.Fprint(buf, "&")
				}
				_, _ = fmt.Fprintf(buf, "%s", key+"="+formValues[key])
			}
			_, _ = fmt.Fprintf(buf, "\" ")
		case MultipartFormType:
			for _, key := range formKeys {
				_, _ = fmt.Fprintf(buf, "-F %s", escape(key+"="+formValues[key]))
			}
			_, _ = fmt.Fprintf(buf, " ")
		}
	}

	_, _ = fmt.Fprint(buf, escape(b.url))

	return buf.String()
}

func (b *CurlBuilder) SetRequest(r *http.Request) *CurlBuilder {
	headerKeys := make([]string, 0, len(r.Header))
	for k := range r.Header {
		headerKeys = append(headerKeys, k)
	}
	sort.Strings(headerKeys)

	headers := make([]string, 0, len(headerKeys)*2)
	for _, headerName := range headerKeys {
		headers = append(headers, headerName, r.Header.Get(headerName))
	}

	if r.MultipartForm == nil {
		formKeys := make([]string, 0, len(r.PostForm))
		for k := range r.Form {
			formKeys = append(formKeys, k)
		}
		sort.Strings(formKeys)

		formValues := make([]string, 0, len(formKeys)*2)
		for _, formName := range formKeys {
			formValues = append(formValues, formName, r.Form.Get(formName))
		}
		b.SetFormValues(MultipartFormType, formValues...)
	} else {
		formKeys := make([]string, 0, len(r.MultipartForm.Value))
		for k := range r.Form {
			formKeys = append(formKeys, k)
		}
		sort.Strings(formKeys)

		formValues := make([]string, 0, len(formKeys)*2)
		for _, formName := range formKeys {
			formValues = append(formValues, formName, r.Form.Get(formName))
		}
		b.SetFormValues(UrlencodedFormType, formValues...)
	}

	var saveBody io.ReadCloser
	saveBody, r.Body, _ = drainBody(r.Body)

	schema := r.URL.Scheme
	requestURL := r.URL.String()

	if schema == "" {
		schema = "http"
		if r.TLS != nil {
			schema = "https"
		}
		requestURL = schema + "://" + r.Host + r.URL.Path
	}

	b.SetBody(saveBody).SetHeaders(headers...).SetMethod(r.Method).SetURL(requestURL)

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

func (b *CurlBuilder) SetFormValues(formType FormType, formValues ...string) *CurlBuilder {
	if len(formValues)%2 != 0 {
		panic("SetFormValues: form values must be key/value")
	}
	b.formType = formType
	b.formValues = formValues
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
