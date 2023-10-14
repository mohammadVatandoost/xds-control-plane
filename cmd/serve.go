package main

import (
	"github.com/mohammadVatandoost/xds-conrol-plane/internal/xds"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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
	conf := loadConfigOrPanic(cmd)
	// configureLoggerOrPanic(conf.Logger)

	log := logger.WithName("main")
	log.Infof("XDS control plane config, ADSEnabled: %v, ListenPort: %v", conf.XDS.ADSEnabled, conf.XDS.ListenPort)
	xdsServer := xds.NewControlPlane(&conf.XDS, nil)
	err := xdsServer.Run()
	if err != nil {
		log.Error(err)
	}
	return nil
}
