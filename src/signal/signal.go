// Package signal registers OS signal handlers for graceful shutdown of cvedex.
package signal

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// Shutdowner is implemented by any service that can be gracefully stopped.
type Shutdowner interface {
	Shutdown()
}

// Handle starts a goroutine that listens for SIGTERM, SIGINT, and SIGHUP
// and calls svc.Shutdown() on receipt.
func Handle(svc Shutdowner) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	go func() {
		sig := <-ch
		slog.Info("signal received, shutting down", "signal", sig)
		svc.Shutdown()
	}()
}
