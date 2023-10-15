package rest

import (
	"errors"
	"fmt"

	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/config"
)



var _ config.Config = &RestAPIConfig{}

type RestAPIConfig struct {
	Port uint32 `json:"port" envconfig:"REST_API_HTTP_PORT"`
	Host string `json:"host" envconfig:"REST_API_HTTP_HOST"`
}

func (c *RestAPIConfig) Validate() error {
	var errs error
	if c.Port > 65535 {
		errs = errors.Join(errs, errors.New("port must be in range [0 65535]"))
	}
	return errs
}

func (c *RestAPIConfig) Sanitize() {
}

func (c *RestAPIConfig) String() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func DefaultRestApiConfig() *RestAPIConfig {
	return &RestAPIConfig{
		Port: 8765,
		Host: "",
	}
}