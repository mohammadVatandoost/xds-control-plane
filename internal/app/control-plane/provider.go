package controlplane

import (
	"sync"

	"github.com/mohammadVatandoost/xds-conrol-plane/internal/node"
	"github.com/mohammadVatandoost/xds-conrol-plane/internal/resource"
	"github.com/mohammadVatandoost/xds-conrol-plane/internal/xds"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/app/controlplane"
)

type App struct {
	conf          *controlplane.ControlPlaneConfig
	nodes         map[string]*node.Node
	mu            sync.RWMutex
	resources     map[string]*resource.Resource
	muResource    sync.RWMutex
	snapshotCache xds.SnapshotCache
}

func NewApp(conf *controlplane.ControlPlaneConfig, snapshotCache xds.SnapshotCache) *App {
	return &App{
		conf:          conf,
		nodes:         make(map[string]*node.Node),
		mu:            sync.RWMutex{},
		resources:     make(map[string]*resource.Resource),
		muResource:    sync.RWMutex{},
		snapshotCache: snapshotCache,
	}
}
