# server

`server` runs related graceful HTTP listeners together. If one listener exits, the others
are cancelled and shut down.

```go
err := server.Run(ctx,
    server.Listener{
        Address: ":8080",
        Handler: application,
        Logger:  logger.With("server", "application"),
    },
    server.Listener{
        Address: ":9090",
        Handler: management,
        Logger:  logger.With("server", "management"),
    },
)
```

Configure the logger before constructing each `Listener`; `server` forwards it unchanged.
Use `Listener.Options` with `grace.With...` options when a listener needs different timeouts
or limits.
