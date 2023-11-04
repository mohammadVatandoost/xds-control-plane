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

// func (cp *ControlPlane) UpdateCache(nodeID string, resourceNames []string, resourceType string) {

// slog.Info("UpdateCache", "nodeID", nodeID, "resourceNames", resourceNames)
// clusters := make([]types.Resource, 0)
// listeners := make([]types.Resource, 0)
// endpoints := make([]types.Resource, 0)
// routes := make([]types.Resource, 0)
// 	resourceNamesMap := make(map[string]struct{}, 0)
// 	for _, v := range resourceNames {
// 		resourceNamesMap[v] = struct{}{}
// 	}

// 	for _, inform := range cp.serviceInformers {
// 		for _, svc := range inform.GetStore().List() {
// 			if reflect.TypeOf(svc).Elem().Name() == "Endpoints" {
// 				continue
// 			}
// 			k8sService, ok := svc.(*corev1.Service)
// 			if !ok {
// 				slog.Error("service type is not match", "type", reflect.TypeOf(svc).Elem().Name())
// 				continue
// 			}
// 			// cp.log.Infof("k8sService Name: %v", k8sService.Name)
// 			_, ok = resourceNamesMap[k8sService.Name]
// 			if !ok {
// 				continue
// 			}
// 			seviceConfig := ServiceConfig{}
// 			// seviceConfig.GRPCServiceName = k8sService.Name
// 			seviceConfig.ServiceName = k8sService.Name
// 			seviceConfig.Namespace = k8sService.Namespace
// 			for _, port := range k8sService.Spec.Ports {
// 				if strings.Contains(port.Name, "grpc") {
// 					seviceConfig.PortName = port.Name
// 					seviceConfig.Protocol = "tcp"
// 					seviceConfig.Region = "us-central1"
// 					seviceConfig.Zone = "us-central1-a"
// 					edsService, clsService, rdsService, lsnrService, err := cp.makeXDSConfigFromService(seviceConfig)
// 					if err != nil {
// 						slog.Error("couldn't make service", "error", err)
// 					}
// 					endpoints = append(endpoints, edsService)
// 					clusters = append(clusters, clsService)
// 					routes = append(routes, rdsService)
// 					listeners = append(listeners, lsnrService)
// 				}

// 			}
// 		}
// 	}

// atomic.AddInt32(&cp.version, 1)

// snapshot, err := cachev3.NewSnapshot(fmt.Sprint(cp.version), map[resource.Type][]types.Resource{
// 	resource.EndpointType: endpoints,
// 	resource.ClusterType:  clusters,
// 	resource.ListenerType: listeners,
// 	resource.RouteType:    routes,
// })
// 	if err != nil {
// 		slog.Error(">>>>>>>>>>  Error creating snapshot", "error", err)
// 		return
// 	}
// slog.Info("snapshotCache IDs: %v, listeners: %v\n", nodeID, listeners)
// err = cp.snapshotCache.SetSnapshot(context.Background(), nodeID, snapshot)
// if err != nil {
// 	slog.Error("couldn't set snapshot", "error", err)
// }
// }
