# CurlBuilder

:triangular_ruler: Билдер CURL запроса с поддержкой http.Request, echo.Context

Для выполнения обратной операции используйте [mholt/curl-to-go](https://github.com/mholt/curl-to-go).

## Пример

```go
import (
    "http"

    "github.com/555f/curlbuilder"
)

data := bytes.NewBufferString(`{"name":"Don Joe"}`)

b := curlbuilder.New()
b.SetURL("http://www.example.com/path/to/page.html?q=1&bar=foo").
    SetBody(data).
	SetMethod("PUT").
    SetHeaders("Content-Type", "application/json")

fmt.Println(b.String())

// Output: Curl -X PUT -d '{"name":"world"}" -H "Content-Type: application/json" http://www.example.com/path/to/page.html?q=1&bar=foo
```

## Пример http.Request

```go
import (
    "http"

    "github.com/555f/curlbuilder"
)

data := bytes.NewBufferString(`{"name":"Don Joe"}`)
req, _ := http.NewRequest("PUT", "http://www.example.com/path/to/page.html?q=1&bar=foo", data)
req.Header.Set("Content-Type", "application/json")


b := curlbuilder.New()
b.SetRequest(req)

fmt.Println(b.String())

// Output: Curl -X PUT -d '{"name":"world"}" -H "Content-Type: application/json" http://www.example.com/path/to/page.html?q=1&bar=foo
```

## Установка

```bash
go get github.com/555f/curlbuilder
```
