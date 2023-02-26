package main

import (
	"fmt"
	"os"

	"github.com/mohammadVatandoost/xds-conrol-plane/internal/config"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const serviceName = "xds_control_plane"

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func loadConfigOrPanic(cmd *cobra.Command) *config.Config {
	conf, err := config.LoadConfig(cmd)
	if err != nil {
		logrus.WithError(err).Panic("Failed to load configurations")
	}
	return conf
}

func configureLoggerOrPanic(loggerConfig logger.Config) {
	if err := logger.Initialize(&loggerConfig); err != nil {
		logrus.WithError(err).Panic("Failed to configure logger")
	}
}
