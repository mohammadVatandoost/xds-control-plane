package node

import (
	"sync"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
)

func NewNode() *Node {
	return &Node{
		Watching:  map[string]struct{}{},
		Clusters:  make([]types.Resource, 0),
		Listeners: make([]types.Resource, 0),
		Endpoints: make([]types.Resource, 0),
		Routes:    make([]types.Resource, 0),
		Version:   0,
	}
}

type Node struct {
	Watching  map[string]struct{} `json:"Watching"`
	mu        sync.RWMutex
	Clusters  []types.Resource `json:"Clusters"`
	Listeners []types.Resource `json:"Listeners"`
	Endpoints []types.Resource `json:"Endpoints"`
	Routes    []types.Resource `json:"Routes"`
	Version   uint64           `json:"Version"`
}
