package config

import (
	"strings"

	"github.com/mohammadVatandoost/xds-conrol-plane/internal/xds"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/logger"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/prometheus"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config the application's configuration structure
type Config struct {
	ConfigFile string
	Logger     logger.Config
	XDS       xds.Config
	Metric     prometheus.Config
}

// LoadConfig loads the config from a file if specified, otherwise from the environment
func LoadConfig(cmd *cobra.Command) (*Config, error) {
	// Setting defaults for this application

	viper.SetDefault("metric.ListenPort", 9000)

	viper.SetDefault("logger.SentryEnabled", false)
	viper.SetDefault("logger.level", "error")

	viper.SetDefault("xds.ListenPort", 8888)
	viper.SetDefault("xds.ADSEnabled", true)

	// Read Config from ENV
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Read Config from Flags
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Read Config from file
	if configFile, err := cmd.Flags().GetString("config-file"); err == nil && configFile != "" {
		viper.SetConfigFile(configFile)

		if err := viper.ReadInConfig(); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	var config Config

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &config, nil
}