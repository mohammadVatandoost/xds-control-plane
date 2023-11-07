package controlplane

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mohammadVatandoost/xds-conrol-plane/internal/xds"
)

// *** callbacks
func (a *App) NewStreamRequest(id string, resourceNames []string, typeURL string) {
	isUpdateNeeded := false
	node, err := a.GetNode(id)
	if err != nil {
		node = a.CreateNode(id)
		isUpdateNeeded = true
	}
	slog.Info("app NewStreamRequest", "id", id, "resourceNames", resourceNames, "typeURL", typeURL)
	for _, rn := range resourceNames {
		node.AddWatching(rn)
		a.AddResourceWatchToNode(id, rn, typeURL)
	}
	if isUpdateNeeded {
		a.UpdateNodeCache(id)
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
		return
	}
	resources := node.GetWatchings()
	node.ClearResources()
	slog.Info("UpdateCache", "nodeID", nodeID, "version", node.GetVersion())
	for _, rn := range resources {
		resource, ok := a.resources[rn]
		if !ok {
			slog.Error("UpdateCache, resource doesn't exist", "resource", rn, "nodeID", nodeID)
			continue
		}
		// resource.ServiceObj.Spec.Ports[0].Name
		//ToDo: later fix loop through each port name
		endPoint, cluster, route, listner, err := xds.MakeXDSResource(resource, a.conf.Region, a.conf.Zone, "")
		if err != nil {
			slog.Error("UpdateCache, failed to Make XDS Resource", "error", err, "resource", rn, "nodeID", nodeID)
			continue
		}
		slog.Info("UpdateCache, resource", "resource", rn, "nodeID", nodeID, "endPoint", endPoint)
		slog.Info("UpdateCache, resource", "resource", rn, "nodeID", nodeID, "cluster", cluster)
		slog.Info("UpdateCache, resource", "resource", rn, "nodeID", nodeID, "listner", listner)
		slog.Info("UpdateCache, resource", "resource", rn, "nodeID", nodeID, "route", route)
		node.AddCluster(cluster)
		node.AddListener(listner)
		node.AddEndpoint(endPoint)
		node.AddRoute(route)
	}

	cache, err := xds.NewSnapshot(
		fmt.Sprintf("%d", node.GetVersion()),
		node.GetEndpoints(),
		node.GetClusters(),
		node.GetListeners(),
		node.GetRoutes(),
	)
	if err != nil {
		slog.Error("UpdateCache, failed to update cache", "error", err, "nodeID", nodeID)
		return
	}
	a.snapshotCache.UpdateCache(context.Background(), nodeID, cache)
	node.IncreaseVersion()
}
