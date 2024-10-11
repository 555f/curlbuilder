package curlbuilder_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/555f/curlbuilder"
)

func TestFromRequest(t *testing.T) {
	type args struct {
		r func() *http.Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success test",
			args: args{
				r: func() *http.Request {
					r, _ := http.NewRequest(http.MethodGet, "http://test.com", nil)
					r.Header.Add("Content-Type", "text/html")
					return r
				},
			},
			want: "curl -X GET -H 'Content-Type: text/html' 'http://test.com'",
		},
		{
			name: "",
			args: args{
				r: func() *http.Request {
					body := map[string]interface{}{
						"name":  "test",
						"bar":   "foo",
						"slice": []int{1, 2, 3},
					}
					data, _ := json.Marshal(body)
					r, _ := http.NewRequest(http.MethodPost, "http://test.com", bytes.NewBuffer(data))
					return r
				},
			},
			want: "curl -X POST -d '{\"bar\":\"foo\",\"name\":\"test\",\"slice\":[1,2,3]}' 'http://test.com'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := curlbuilder.New()
			b.SetRequest(tt.args.r())
			if got := b.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
