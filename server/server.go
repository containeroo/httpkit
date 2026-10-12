// Package server runs related HTTP listeners with coordinated graceful shutdown.
package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/containeroo/httpkit/grace"
	"golang.org/x/sync/errgroup"
)

// Listener describes one HTTP listener managed by Run.
type Listener struct {
	Address string
	Handler http.Handler
	Logger  *slog.Logger
	Options []grace.Option
}

// Run serves every listener and stops the others when one returns.
func Run(ctx context.Context, listeners ...Listener) error {
	listenerContext, cancel := context.WithCancel(ctx)
	defer cancel()
	group, groupContext := errgroup.WithContext(listenerContext)
	for _, listener := range listeners {
		listener := listener
		group.Go(func() error {
			defer cancel()
			return grace.Run(groupContext, listener.Address, listener.Handler, listener.Logger, listener.Options...)
		})
	}
	return group.Wait()
}
