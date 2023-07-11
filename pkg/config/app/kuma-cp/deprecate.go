package kuma_cp

import (
	"io"

	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/config"
)

var deprecations = []config.Deprecation{}

func PrintDeprecations(cfg *Config, out io.Writer) {
	config.PrintDeprecations(deprecations, cfg, out)
}
