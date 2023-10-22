package node

import (
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
)

func (n *Node) AddWatcher(resource string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.watchers[resource] = struct{}{}
}

func (n *Node) IsWatched(resource string) bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	_, ok := n.watchers[resource]
	return ok
}

func (n *Node) ClearResources() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.clusters = make([]types.Resource, 0)
	n.listeners = make([]types.Resource, 0)
	n.endpoints = make([]types.Resource, 0)
	n.routes = make([]types.Resource, 0)
}

func (n *Node) GetClusters() []types.Resource {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.clusters
}

func (n *Node) AddCluster(r types.Resource) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.clusters = append(n.clusters, r)
}

func (n *Node) GetListeners() []types.Resource {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.listeners
}

func (n *Node) AddListener(r types.Resource) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.listeners = append(n.listeners, r)
}

func (n *Node) GetEndpoints() []types.Resource {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.endpoints
}

func (n *Node) AddEndpoint(r types.Resource) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.endpoints = append(n.endpoints, r)
}

func (n *Node) GetRoutes() []types.Resource {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.routes
}

func (n *Node) AddRoute(r types.Resource) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.routes = append(n.routes, r)
}