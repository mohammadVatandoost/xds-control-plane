package prometheus

import (
	"context"
	"fmt"

	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func StartMetricServerOrPanic(listenPort int) *http.Server {
	prometheusServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", listenPort),
		Handler: promhttp.Handler(),
	}

	go listenAndServeMetrics(prometheusServer)
	return prometheusServer
}

func listenAndServeMetrics(server *http.Server) {
	if err := server.ListenAndServe(); err != nil {
		logrus.Panic(err.Error(), "failed to start liveness http probe listener")
	}
}

func ShutdownMetricServerOrPanic(server *http.Server) {
	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Panic(err, "Failed to shutdown prometheus metric server")
	}
}
