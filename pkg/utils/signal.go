package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func WithSignalCancellation(
	ctx context.Context,
) (context.Context, context.CancelFunc) {
	gracefulStop := make(chan os.Signal, 2)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)

	newContext, cancel := context.WithCancel(ctx)
	go func() {
		<-gracefulStop
		cancel()
	}()
	return newContext, cancel
}
