# parse

`parse` converts HTTP configuration values into types from Go's standard library.

```go
import "github.com/containeroo/httpkit/parse"
```

## Status codes

`ParseStatusCodes` accepts comma-separated HTTP status codes and inclusive ranges. By default, values must be between **100 and 599**. The parser checks range endpoints before expanding them, so extremely large ranges are rejected without allocating large slices.

```go
statusCodes, err := parse.ParseStatusCodes("200,300-302,404")
if err != nil {
    return err
}

// statusCodes is []int{200, 300, 301, 302, 404}
```

To restrict the accepted values, use `WithStatusCodeRange`. Boundaries are inclusive and must be within the HTTP range 100–599, with `min <= max`. Invalid bounds (such as `599, 200`) return an error.

```go
statusCodes, err := parse.ParseStatusCodes(
    "200-204,429,500-503",
    parse.WithStatusCodeRange(200, 599),
)
if err != nil {
    return err
}
```

## Headers

`ParseHeaders` accepts `Key=Value` entries and returns an `http.Header`. Header names are canonicalized, and commas or additional `=` characters in values are preserved.

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
