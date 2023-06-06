package xds

import "sync"

type Node struct {
	watchers map[string]struct{}
	mu sync.RWMutex
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
