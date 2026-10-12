# grace

`grace` runs one `net/http` server and shuts it down cleanly when its context is cancelled.

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()

logger := slog.Default().With("server", "public")
err := grace.Run(ctx, ":8080", handler, logger,
    grace.WithReadHeaderTimeout(10*time.Second),
    grace.WithShutdownTimeout(15*time.Second),
)
```

It supplies conservative defaults for header read, write, idle, and shutdown timeouts.
Use the focused `With...` options to override them. Pass `nil` as the logger to disable
lifecycle logging.
