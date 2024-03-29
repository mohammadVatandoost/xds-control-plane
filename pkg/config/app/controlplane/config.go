package controlplane

import (
	"errors"

	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/config"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/rest"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/xds"
)

var _ config.Config = &ControlPlaneConfig{}

type ControlPlaneConfig struct {
	RestAPIConfig *rest.RestAPIConfig `restAPIConfig:"port" `
	XDSConfig     *xds.XDSConfig      `json:"xdsConfig" `
	Region        string              `json:"region" envconfig:"REGION"`
	Zone          string              `json:"zone" envconfig:"ZONE"`
}

func (c *ControlPlaneConfig) Validate() error {
	var errs error
	if err := c.RestAPIConfig.Validate(); err != nil {
		errs = errors.Join(errs, errors.New("restAPIConfig validation is failed"))
	}
	if err := c.XDSConfig.Validate(); err != nil {
		errs = errors.Join(errs, errors.New("restAPIConfig validation is failed"))
	}
	return errs
}

func (c *ControlPlaneConfig) Sanitize() {
}

func DefaultControlPlaneConfig() *ControlPlaneConfig {
	return &ControlPlaneConfig{
		RestAPIConfig: rest.DefaultRestApiConfig(),
		XDSConfig:     xds.DefaultXDSConfig(),
		Region:        "us-central1",
		Zone:          "us-central1-a",
	}
}
