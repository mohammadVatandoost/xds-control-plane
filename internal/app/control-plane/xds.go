package controlplane

import (
	"log/slog"

	"github.com/mohammadVatandoost/xds-conrol-plane/internal/resource"
	"github.com/mohammadVatandoost/xds-conrol-plane/internal/xds"
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
	resources := node.GetWatchings()
	node.ClearResources()
	slog.Info("UpdateCache", "nodeID", nodeID)
	for _, rn := range resources {
		resource, ok := a.resources[rn]
		if !ok {
			slog.Error("UpdateCache, resource doesn't exist", "resource", rn, "nodeID", nodeID)
			continue
		}
		//ToDo: later fix loop through each port name
		endPoint, cluster, listner, route, err := xds.MakeXDSResource(resource, a.conf.Region, a.conf.Zone, resource.ServiceObj.Spec.Ports[0].Name)
		if err != nil {
			slog.Error("UpdateCache, failed to Make XDS Resource", "error", err, "resource", rn, "nodeID", nodeID)
			continue
		}
		node.AddCluster(cluster)
		node.AddListener(listner)
		node.AddEndpoint(endPoint)
		node.AddRoute(route)
	}
}