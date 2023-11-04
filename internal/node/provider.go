package node

import (
	"sync"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
)

func NewNode() *Node {
	return &Node{
		watching:  map[string]struct{}{},
		clusters:  make([]types.Resource, 0),
		listeners: make([]types.Resource, 0),
		endpoints: make([]types.Resource, 0),
		routes:    make([]types.Resource, 0),
		version:   0,
	}
}

type Node struct {
	watching  map[string]struct{}
	mu        sync.RWMutex
	clusters  []types.Resource
	listeners []types.Resource
	endpoints []types.Resource
	routes    []types.Resource
	version   uint64
}
