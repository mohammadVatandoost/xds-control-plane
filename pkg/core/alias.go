package core

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/logger"
	// kube_log "sigs.k8s.io/controller-runtime/pkg/log"
	// kuma_log "github.com/mohammadVatandoost/xds-conrol-plane/pkg/log"
)

var (
	// TODO remove dependency on kubernetes see: https://github.com/mohammadVatandoost/xds-conrol-plane/issues/2798
	// Log                   = kube_log.Log
	// NewLogger             = kuma_log.NewLogger
	// NewLoggerTo           = kuma_log.NewLoggerTo
	// NewLoggerWithRotation = kuma_log.NewLoggerWithRotation
	// SetLogger             = kube_log.SetLogger
	log     = logger.WithName("pkg/core")
	Now     = time.Now
	TempDir = os.TempDir

	SetupSignalHandler = func() (context.Context, context.Context) {
		gracefulCtx, gracefulCancel := context.WithCancel(context.Background())
		ctx, cancel := context.WithCancel(context.Background())
		c := make(chan os.Signal, 3)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			s := <-c
			log.Info("Received signal, stopping instance gracefully", "signal", s.String())
			gracefulCancel()
			s = <-c
			log.Info("Received second signal, stopping instance", "signal", s.String())
			cancel()
			s = <-c
			log.Info("Received third signal, force exit", "signal", s.String())
			os.Exit(1)
		}()
		return gracefulCtx, ctx
	}
)

func NewUUID() string {
	return uuid.NewString()
}
