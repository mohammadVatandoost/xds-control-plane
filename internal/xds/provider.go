package xds

import (
	"context"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	xds "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/mohammadVatandoost/xds-conrol-plane/internal/node"
	xdsConfig "github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/xds"
)

func NewControlPlane(config *xdsConfig.XDSConfig) *ControlPlane {
	snapshotCache := cache.NewSnapshotCache(config.ADSEnabled, cache.IDHash{}, nil)

	cp := &ControlPlane{
		version:       0,
		snapshotCache: snapshotCache,
		conf:          config,
		nodes:         make(map[string]*node.Node),
		resources:     make(map[string]map[string]struct{}),
	}
	// callBacks := newCallBack(log, cp)
	// cp.callBacks = callBacks
	cp.server = xds.NewServer(context.Background(), snapshotCache, cp)
	return cp
}

