# middleware

`middleware` provides standard-library `http.Handler` wrappers. It deliberately does
not choose an authentication model, response format, or logger schema.

```go
handler := middleware.Chain(app,
    middleware.SecurityHeaders(http.Header{
        "X-Content-Type-Options": {"nosniff"},
    }),
    middleware.AccessLog(logger),
    middleware.RecoverPanics(writePanicResponse, responseCommitted),
    middleware.RejectCrossOriginWrites(logger, writeForbidden),
)
```

Middleware executes in declaration order. `RecoverPanics` delegates the error response
to your application; `RejectCrossOriginWrites` delegates rejected-request responses too.
This preserves each service's response contract.
