package curlbuilder_test

import (
	"testing"

	"github.com/555f/curlbuilder"
)

type TestBody struct {
	Name string `json:"name"`
}

func TestCurlBuilder_String(t *testing.T) {
	type fields struct {
		url     string
		body    interface{}
		method  string
		headers []string
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
			want: "curl -X POST -d '{\"bar\":\"foo\",\"name\":\"test\",\"slice\":[1,2,3]}' http://test.com",
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
			want: "curl -k -X POST -d '{\"name\":\"test\"}' https://test.com",
		},
		{
			name: "Test 3",
			fields: fields{
				url:     "http://test.com",
				method:  "POST",
				headers: nil,
			},
			want: "curl -X POST http://test.com",
		},
		{
			name: "Test 4",
			fields: fields{
				url:     "http://test.com",
				method:  "GET",
				headers: nil,
			},
			want: "curl -X GET http://test.com",
		},
		{
			name: "Test 5s",
			fields: fields{
				url:     "http://test.com",
				method:  "GET",
				headers: []string{"Content-Type", "text/html"},
			},
			want: "curl -X GET -H 'Content-Type: text/html' http://test.com",
		},
		{
			name: "Test 6",
			fields: fields{
				url:     "http://test.com",
				method:  "GET",
				headers: []string{"Content-Type", "text/html", "Accept", "application/xml"},
			},
			want: "curl -X GET -H 'Accept: application/xml' -H 'Content-Type: text/html' http://test.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := curlbuilder.New()
			b.SetURL(tt.fields.url).
				SetBody(tt.fields.body).
				SetMethod(tt.fields.method).
				SetHeaders(tt.fields.headers...)
			if got := b.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
