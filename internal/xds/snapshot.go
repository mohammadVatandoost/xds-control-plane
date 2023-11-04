package xds

import (
	"context"
	"fmt"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
)

type SnapshotCache interface {
	UpdateCache(ctx context.Context, nodeID string, snapshot *cachev3.Snapshot) error
}

func NewSnapshotCache(ADSEnabled bool) *XDSSnapshotCache {
	return &XDSSnapshotCache{
		cachev3.NewSnapshotCache(ADSEnabled, cachev3.IDHash{}, nil),
	}
}

type XDSSnapshotCache struct {
	cachev3.SnapshotCache
}

func (xc *XDSSnapshotCache) UpdateCache(ctx context.Context, nodeID string, snapshot *cachev3.Snapshot) error {
	return xc.SetSnapshot(context.Background(), nodeID, snapshot)
}

func NewSnapshot(version string, endpoints, clusters, listeners, routes []types.Resource) (*cachev3.Snapshot, error) {
	if version == "" {
		return nil, fmt.Errorf("version is empty")
	}
	if endpoints == nil {
		endpoints = []types.Resource{}
	}
	if clusters == nil {
		clusters = []types.Resource{}
	}
	if listeners == nil {
		listeners = []types.Resource{}
	}
	if routes == nil {
		routes = []types.Resource{}
	}
	return cachev3.NewSnapshot(version, map[string][]types.Resource{
		resource.EndpointType: endpoints,
		resource.ClusterType:  clusters,
		resource.ListenerType: listeners,
		resource.RouteType:    routes,
	})
}
