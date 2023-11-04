package xds

import (
	"context"

	xds "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/mohammadVatandoost/xds-conrol-plane/internal/node"
	xdsConfig "github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/xds"
)

func NewControlPlane(config *xdsConfig.XDSConfig, app App, cache *XDSSnapshotCache) *ControlPlane {
	cp := &ControlPlane{
		version:   0,
		cache:     cache,
		conf:      config,
		nodes:     make(map[string]*node.Node),
		resources: make(map[string]map[string]struct{}),
		app:       app,
	}
	cp.server = xds.NewServer(context.Background(), cache, cp)
	return cp
}
