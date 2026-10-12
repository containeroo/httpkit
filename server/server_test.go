package server

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRunStopsPeerOnBindFailure(t *testing.T) {
	t.Parallel()
	occupied, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = occupied.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = Run(ctx,
		Listener{Address: occupied.Addr().String(), Handler: http.NotFoundHandler(), Logger: slog.New(slog.NewTextHandler(io.Discard, nil))},
		Listener{Address: "127.0.0.1:0", Handler: http.NotFoundHandler(), Logger: slog.New(slog.NewTextHandler(io.Discard, nil))},
	)
	require.Error(t, err)
	require.NoError(t, ctx.Err(), "listener failure should stop the peer before the deadline")
}
