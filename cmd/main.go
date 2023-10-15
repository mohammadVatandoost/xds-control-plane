package main

import (
	"os"
	"log/slog"

	"github.com/mohammadVatandoost/xds-conrol-plane/internal/xds"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/config"
	controlplaneConfig "github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/app/controlplane"
)

const serviceName = "xds_control_plane"

func main() {
	exitCode := 0
	defer func ()  {
		os.Exit(exitCode)
	}()
	// conf := loadConfigOrPanic(cmd)
	// configureLoggerOrPanic(conf.Logger)

	conf := controlplaneConfig.DefaultControlPlaneConfig()
	err := config.Load("", conf)
	if err != nil {
		slog.Error("couldn't load configs", "error", err)
		exitCode = -1
		return
	}

	slog.Info("XDS control plane config", "XDS.ADSEnabled", conf.XDSConfig.ADSEnabled, "ListenPort", conf.XDSConfig.Port)
	xdsServer := xds.NewControlPlane(conf.XDSConfig)
	err = xdsServer.Run()
	if err != nil {
		slog.Error("couldn't run xdsServer", "error", err)
	}
}

