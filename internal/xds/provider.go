package xds

import (
	"context"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	xds "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/sirupsen/logrus"
)

func NewControlPlane(log *logrus.Logger, config *Config, storage cache.Storage) *ControlPlane {
	snapshotCache := cache.NewSnapshotCache(config.ADSEnabled, cache.IDHash{}, log)
	if storage != nil {
		snapshotCache = cache.NewSnapshotCacheWithStorage(config.ADSEnabled, cache.IDHash{}, log, storage)
	}
	cp := &ControlPlane{
		log:           log,
		version:       0,
		snapshotCache: snapshotCache,
		storage:       storage,
		conf:          config,
	}
	callBacks := newCallBack(log, cp)
	cp.callBacks = callBacks
	cp.server = xds.NewServer(context.Background(), snapshotCache, callBacks)
	return cp
}

func newCallBack(log *logrus.Logger, eventsHandler EventsHandler) *callbacks {
	signal := make(chan struct{})
	return &callbacks{
		log:           log,
		signal:        signal,
		fetches:       0,
		requests:      0,
		eventsHandler: eventsHandler,
	}
}
