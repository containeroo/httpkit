# prefix

`prefix` lets the same application run at `/` locally and below a reverse-proxy path in production.

```go
prefixValue := prefix.NormalizeRoutePrefix("/wiki/") // "/wiki"
handler := prefix.MountUnderPrefixWithOptions(app, prefixValue,
    prefix.WithRedirectRewriting(),
)

pageURL := prefix.RouteURL(prefixValue, "/pages/start") // "/wiki/pages/start"
```

Application routes should stay prefix-free. Use `RouteURL` for generated links and
`URLForRequest` or `Redirect` for request-relative links and redirects.
