package controlplane

import (
	"log/slog"

	"github.com/mohammadVatandoost/xds-conrol-plane/internal/resource"
)

// *** callbacks
func (a *App) NewStreamRequest(id string, resourceNames []string, typeURL string) {
	node := a.CreateNode(id)
	slog.Info("app NewStreamRequest", "id", id, "resourceNames", resourceNames, "typeURL", typeURL)
	for _, rn := range resourceNames {
		node.AddWatching(rn)
		a.AddResourceWatchToNode(id, rn, typeURL)
	}
}

func (a *App) StreamClosed(id string) {
	err := a.DeleteNode(id)
	if err != nil {
		slog.Error("app stream closed", "error", err, "id", id)
	}
}
// ***

func (a *App) UpdateNodeCache(nodeID string) {
	node, ok := a.nodes[nodeID]
	if !ok {
		slog.Error("UpdateNodeCache, node doesn't exist", "nodeID", nodeID)
	}
	resources := node.
	slog.Info("UpdateCache", "nodeID", nodeID)
	clusters := make([]types.Resource, 0)
	listeners := make([]types.Resource, 0)
	endpoints := make([]types.Resource, 0)
	routes := make([]types.Resource, 0)
}