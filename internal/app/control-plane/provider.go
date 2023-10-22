package controlplane

import (
	"sync"

	"github.com/mohammadVatandoost/xds-conrol-plane/internal/node"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/app/controlplane"
)

type App struct {
	conf *controlplane.ControlPlaneConfig
	nodes             map[string]*node.Node
	mu                sync.RWMutex
	resources         map[string]map[string]struct{} // A resource is watched by which nodes
	muResource        sync.RWMutex
}


func NewApp(conf *controlplane.ControlPlaneConfig) *App {
	return &App{
		conf: conf,
	}
}