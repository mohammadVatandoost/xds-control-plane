package xds

import (
	"sync"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
)

func NewNode() *Node {
	return &Node{
		watchers:  map[string]struct{}{},
		clusters:  make([]types.Resource, 0),
		listeners: make([]types.Resource, 0),
		endpoints: make([]types.Resource, 0),
		routes:    make([]types.Resource, 0),
	}
}

type Node struct {
	watchers  map[string]struct{}
	mu        sync.RWMutex
	clusters  []types.Resource
	listeners []types.Resource
	endpoints []types.Resource
	routes    []types.Resource
}

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
