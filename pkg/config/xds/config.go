package xds


import (
	"errors"

	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/config"
	"go.uber.org/multierr"
)



var _ config.Config = &XDSConfig{}

type XDSConfig struct {
	Port uint32 `json:"port" envconfig:"XDS_PORT"`
	ADSEnabled bool `json:"adsEnable" envconfig:"ADS_ENABLE"`
}

func (c *XDSConfig) Validate() error {
	var errs error
	if c.Port > 65535 {
		errs = multierr.Append(errs, errors.New("port must be in range [0 65535]"))
	}
	return errs
}

func (c *XDSConfig) Sanitize() {
}

func DefaultXDSConfig() *XDSConfig {
	return &XDSConfig{
		Port: 8765,
		ADSEnabled: true,
	}
}