package main

import (
	"github.com/mohammadVatandoost/xds-conrol-plane/internal/xds"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/config"
	controlplaneConfig "github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/app/controlplane"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var log = logger.NewLoggerWithName("main")

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := serve(cmd, args); err != nil {
			logrus.WithError(err).Fatal("Failed to serve.")
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) error {
	printVersion()
	conf := controlplaneConfig.DefaultControlPlaneConfig()
	err := config.Load("", conf)
	if err != nil {
		log.Error(err, "couldn't load config")
		return err
	}
	log.Info("XDS control plane config", "ADSEnabled", conf.XDSConfig.ADSEnabled, "ListenPort", conf.XDSConfig.Port)
	xdsServer := xds.NewControlPlane(conf.XDSConfig, nil)
	err = xdsServer.Run()
	if err != nil {
		log.Error(err, "can not run xds server")
	}
	return nil
}
