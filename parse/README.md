# parse

`parse` converts HTTP configuration values into types from Go's standard library.

```go
import "github.com/containeroo/httpkit/parse"
```

## Status codes

`ParseStatusCodes` accepts comma-separated status codes and inclusive ranges.

```go
statusCodes, err := parse.ParseStatusCodes("200,300-302,404")
if err != nil {
	return err
}

// statusCodes is []int{200, 300, 301, 302, 404}
```

## Headers

`ParseHeaders` accepts `Key=Value` entries and returns an `http.Header`. Header
names are canonicalized, and commas or additional `=` characters in values are
preserved.

```go
headers, err := parse.ParseHeaders([]string{
	"Content-Type=application/json",
	"Accept=text/html, application/json",
	"X-Trace=one",
	"X-Trace=two",
}, true)
if err != nil {
	return err
}

// headers.Values("X-Trace") is []string{"one", "two"}
```

Pass `false` as the second argument to reject duplicate header names.
