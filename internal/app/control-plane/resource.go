package controlplane

import (
	"fmt"

	"github.com/mohammadVatandoost/xds-conrol-plane/internal/node"
)

func (a *App) GetServices() {
}

func (a *App) CreateNode(id string) *node.Node {
	a.mu.Lock()
	defer a.mu.Unlock()
	n, ok := a.nodes[id]
	if !ok {
		n = node.NewNode()
	}
	a.nodes[id] = n
	return n
}

func (a *App) GetNode(id string) (*node.Node, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	node, ok := a.nodes[id]
	if !ok {
		return nil, fmt.Errorf("node with id: %s is not exist", id)
	}
	return node, nil
}

func (a *App) DeleteNode(id string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	_, ok := a.nodes[id]
	if !ok {
		return fmt.Errorf("node with id: %s is not exist", id)
	}
	delete(a.nodes, id)
	return nil
}

func (a *App) AddResourceWatchToNode(id string, resource string) {
	a.muResource.Lock()
	defer a.muResource.Unlock()
	nodes, ok := a.resources[resource]
	if !ok {
		nodes = make(map[string]struct{})
		a.resources[resource] = nodes
	}
	nodes[id] = struct{}{}
}

func (a *App) GetNodesWatchTheResource(resource string) []string {
	a.muResource.RLock()
	defer a.muResource.RUnlock()
	nodesArray := make([]string, 0)
	nodes, ok := a.resources[resource]
	if !ok {
		return nodesArray
	}
	for n := range nodes {
		nodesArray = append(nodesArray, n)
	}
	return nodesArray
}
