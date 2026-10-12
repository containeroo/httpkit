# httpkit

Small, composable HTTP building blocks for Go's standard library.

| Package                      | Purpose                                                                                  |
| ---------------------------- | ---------------------------------------------------------------------------------------- |
| [`grace`](./grace)           | Run one HTTP server with bounded graceful shutdown.                                      |
| [`middleware`](./middleware) | Compose handlers, log requests, recover panics, and enforce browser-origin write checks. |
| [`prefix`](./prefix)         | Mount handlers and generate URLs below an optional path prefix.                          |
| [`response`](./response)     | Write JSON/XML responses, track final status, and validate local redirects.              |
| [`server`](./server)         | Run related graceful listeners as one unit.                                              |

## Install

```bash
go get github.com/containeroo/httpkit
```

Packages are independent: import only what a service needs. Logging, error responses,
authentication, and application policy remain owned by the consuming service.

`grace` and `prefix` began as additive extractions from `httpgrace` and
`httpprefix`; the original modules remain available while consumers migrate.
