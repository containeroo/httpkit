# httpkit

Small, composable HTTP building blocks for Go's standard library.

| Package                      | Purpose                                                                                  |
| ---------------------------- | ---------------------------------------------------------------------------------------- |
| [`assets`](./assets)         | Serve embedded files with asset caching and gzip negotiation.                            |
| [`grace`](./grace)           | Run one HTTP server with bounded graceful shutdown.                                      |
| [`middleware`](./middleware) | Compose handlers, log requests, recover panics, and enforce browser-origin write checks. |
| [`parse`](./parse)           | Parse header entries and status-code expressions from configuration.                     |
| [`prefix`](./prefix)         | Mount handlers and generate URLs below an optional path prefix.                          |
| [`response`](./response)     | Write JSON/XML responses, track final status, and validate local redirects.              |
| [`server`](./server)         | Run related graceful listeners as one unit.                                              |

## Install

```bash
go get github.com/containeroo/httpkit
```

`assets.Static` serves embedded frontend files with cache control and optional gzip.
`response.NewRequestWriter` and `response.LogFailure` correlate HTTP failures with access logs.
