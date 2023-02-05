package main

import (
	"xds-conrol-plane/internal/xds"
	"xds-conrol-plane/pkg/logger"
)

func main() {
	logger.Initialize(&logger.Config{Level: "debug"})
	log := logger.NewLogger()
	xdsServer := xds.NewControlPlane(log)
	err := xdsServer.Run()
	if err != nil {
		log.Error(err)
	}
}
