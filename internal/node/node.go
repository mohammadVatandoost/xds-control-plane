package node

import (
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
)

func (n *Node) AddWatching(resource string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Watching[resource] = struct{}{}
}

func (n *Node) IsWatched(resource string) bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	_, ok := n.Watching[resource]
	return ok
}

func (n *Node) GetWatchings() []string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	Watchings := make([]string, 0)
	for w := range n.Watching {
		Watchings = append(Watchings, w)
	}
	return Watchings
}

func (n *Node) ClearResources() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Clusters = make([]types.Resource, 0)
	n.Listeners = make([]types.Resource, 0)
	n.Endpoints = make([]types.Resource, 0)
	n.Routes = make([]types.Resource, 0)
}

func (n *Node) GetClusters() []types.Resource {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Clusters
}

func (n *Node) AddCluster(r types.Resource) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Clusters = append(n.Clusters, r)
}

func (n *Node) GetListeners() []types.Resource {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Listeners
}

func (n *Node) AddListener(r types.Resource) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Listeners = append(n.Listeners, r)
}

func (n *Node) GetEndpoints() []types.Resource {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Endpoints
}

func (n *Node) AddEndpoint(r types.Resource) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Endpoints = append(n.Endpoints, r)
}

func (n *Node) GetRoutes() []types.Resource {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Routes
}

func (n *Node) AddRoute(r types.Resource) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Routes = append(n.Routes, r)
}

func (n *Node) GetVersion() uint64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.Version
}

func (n *Node) IncreaseVersion() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Version++
}
