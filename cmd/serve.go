package main

import (
	"github.com/mohammadVatandoost/xds-conrol-plane/internal/xds"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/logger"
	api_server_config "github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/api-server"
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
    cfg := api_server_config.DefaultConfig()
	conf := loadConfigOrPanic(cmd)
	// configureLoggerOrPanic(conf.Logger)
	apiServer := api_server.NewApiServer()

	log := logger.WithName("main")
	log.Infof("XDS control plane config, ADSEnabled: %v, ListenPort: %v", conf.XDS.ADSEnabled, conf.XDS.ListenPort)
	xdsServer := xds.NewControlPlane(&conf.XDS, nil)
	err := xdsServer.Run()
	if err != nil {
		log.Error(err)
	}
	return nil
}
