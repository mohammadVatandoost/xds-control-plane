package rest

import (
	"errors"

	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/config"
	"go.uber.org/multierr"
)



var _ config.Config = &RestAPIConfig{}

type RestAPIConfig struct {
	Port uint32 `json:"port" envconfig:"REST_API_HTTP_PORT"`
}

func (c *RestAPIConfig) Validate() error {
	var errs error
	if c.Port > 65535 {
		errs = multierr.Append(errs, errors.New("port must be in range [0 65535]"))
	}
	return errs
}

func (c *RestAPIConfig) Sanitize() {
}

func DefaultRestApiConfig() *RestAPIConfig {
	return &RestAPIConfig{
		Port: 8765,
	}
}