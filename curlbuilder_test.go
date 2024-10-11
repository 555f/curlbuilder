package curlbuilder_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/555f/curlbuilder"
)

type TestBody struct {
	Name string `json:"name"`
}

func TestCurlBuilder_String(t *testing.T) {
	type fields struct {
		url        string
		body       interface{}
		method     string
		formType   curlbuilder.FormType
		formValues []string
		headers    []string
		secrets    []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test 1",
			fields: fields{
				url: "http://test.com",
				body: map[string]interface{}{
					"name":  "test",
					"bar":   "foo",
					"slice": []int{1, 2, 3},
				},
				method:  "POST",
				headers: nil,
			},
			want: "curl -X POST -d '{\"bar\":\"foo\",\"name\":\"test\",\"slice\":[1,2,3]}' 'http://test.com'",
		},
		{
			name: "Test 2",
			fields: fields{
				url: "https://test.com",
				body: map[string]interface{}{
					"name": "test",
				},
				method:  "POST",
				headers: nil,
			},
			want: "curl -k -X POST -d '{\"name\":\"test\"}' 'https://test.com'",
		},
		{
			name: "Test 3",
			fields: fields{
				url:     "http://test.com",
				method:  "POST",
				headers: nil,
			},
			want: "curl -X POST 'http://test.com'",
		},
		{
			name: "Test 4",
			fields: fields{
				url:     "http://test.com",
				method:  "GET",
				headers: nil,
			},
			want: "curl -X GET 'http://test.com'",
		},
		{
			name: "Test 5s",
			fields: fields{
				url:     "http://test.com",
				method:  "GET",
				headers: []string{"Content-Type", "text/html"},
			},
			want: "curl -X GET -H 'Content-Type: text/html' 'http://test.com'",
		},
		{
			name: "Test 6",
			fields: fields{
				url:     "http://test.com",
				method:  "GET",
				headers: []string{"Content-Type", "text/html", "Accept", "application/xml"},
			},
			want: "curl -X GET -H 'Accept: application/xml' -H 'Content-Type: text/html' 'http://test.com'",
		},
		{
			name: "Test 7",
			fields: fields{
				url:        "http://test.com",
				method:     "POST",
				formType:   curlbuilder.UrlencodedFormType,
				formValues: []string{"key", "val"},
			},
			want: "curl -X POST -d \"key=val\" 'http://test.com'",
		},
		{
			name: "Test 8",
			fields: fields{
				url:        "http://test.com",
				method:     "POST",
				formType:   curlbuilder.MultipartFormType,
				formValues: []string{"key", "val"},
			},
			want: "curl -X POST -F 'key=val' 'http://test.com'",
		},
		{
			name: "Test 9",
			fields: fields{
				url:     "https://test.com",
				body:    "test",
				method:  "POST",
				headers: nil,
			},
			want: "curl -k -X POST -d 'test' 'https://test.com'",
		},
		{
			name: "Test 10",
			fields: fields{
				url:     "https://test.com",
				body:    []byte("test"),
				method:  "POST",
				headers: nil,
			},
			want: "curl -k -X POST -d 'test' 'https://test.com'",
		},
		{
			name: "Test 11",
			fields: fields{
				url:     "https://test.com",
				body:    []byte("test"),
				method:  "POST",
				headers: []string{"Authorization", "Bearer 123123"},
				secrets: []string{"Authorization"},
			},
			want: "curl -k -X POST -d 'test' -H 'Authorization: *****' 'https://test.com'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := curlbuilder.New()
			if tt.fields.formType != "" && len(tt.fields.formValues) > 0 {
				b.SetFormValues(tt.fields.formType, tt.fields.formValues...)
			}
			b.SetURL(tt.fields.url).
				SetBody(tt.fields.body).
				SetMethod(tt.fields.method).
				SetHeaders(tt.fields.headers...).
				SetSecret(tt.fields.secrets...)
			if got := b.String(); got != tt.want {
				t.Errorf("\n got  = %v\n want = %v\n\n", got, tt.want)
			}
		})
	}
}

func BenchmarkGetCurlCommand(b *testing.B) {
	form := url.Values{}

	for i := 0; i <= b.N; i++ {
		form.Add("number", strconv.Itoa(i))
		body := form.Encode()
		req, _ := http.NewRequest(http.MethodPost, "http://foo.com", io.NopCloser(bytes.NewBufferString(body)))
		curlbuilder.FromRequest(req)
	}
}

func TestGetCurlCommand_serverSide(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := curlbuilder.FromRequest(r)
		fmt.Fprint(w, c)
	}))
	defer svr.Close()

	resp, err := http.Get(svr.URL)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	exp := fmt.Sprintf("curl -X GET -H 'Accept-Encoding: gzip' -H 'User-Agent: Go-http-client/1.1' '%s/'", svr.URL)
	if out := string(data); out != exp {
		t.Errorf("act: %s, exp: %s", out, exp)
	}
}
