package xds

import (
	"context"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	xds "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/logger"
)

var log = logger.NewLoggerWithName("internal/xds")

func NewControlPlane(config *Config, storage cache.Storage) *ControlPlane {

	snapshotCache := cache.NewSnapshotCache(config.ADSEnabled, cache.IDHash{}, nil)
	if storage != nil {
		snapshotCache = cache.NewSnapshotCacheWithStorage(config.ADSEnabled, cache.IDHash{}, nil, storage)
	}
	cp := &ControlPlane{
		version:       0,
		snapshotCache: snapshotCache,
		storage:       storage,
		conf:          config,
		nodes:         make(map[string]*Node),
		resources:     make(map[string]map[string]struct{}),
	}
	// callBacks := newCallBack(log, cp)
	// cp.callBacks = callBacks
	cp.server = xds.NewServer(context.Background(), snapshotCache, cp)
	return cp
}

// func newCallBack(log *logrus.Logger, eventsHandler EventsHandler) *callbacks {
// 	signal := make(chan struct{})
// 	return &callbacks{
// 		log:           log,
// 		signal:        signal,
// 		fetches:       0,
// 		requests:      0,
// 		eventsHandler: eventsHandler,
// 	}
// }
