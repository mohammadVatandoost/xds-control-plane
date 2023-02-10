package xds

import (
	"context"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	xds "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/sirupsen/logrus"
)

func NewControlPlane(log *logrus.Logger) *ControlPlane {
	snapshotCache := cache.NewSnapshotCache(true, cache.IDHash{}, log)
	snapshotCache.GetStatusKeys()
	callBacks := newCallBack(log)
	return &ControlPlane{
		log:           log,
		version:       0,
		snapshotCache: snapshotCache,
		callBacks:     callBacks,
		server:        xds.NewServer(context.Background(), snapshotCache, callBacks),
	}
}

func newCallBack(log *logrus.Logger) *callbacks {
	signal := make(chan struct{})
	return &callbacks{
		log:      log,
		signal:   signal,
		fetches:  0,
		requests: 0,
	}
}
