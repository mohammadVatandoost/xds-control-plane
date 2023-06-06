package xds

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync/atomic"

	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	corev1 "k8s.io/api/core/v1"
)

func (cp *ControlPlane) HandleServicesUpdate(oldObj, newObj interface{}) {
	// cp.log.Info("ControlPlane HandleServicesUpdate")
	clusters := make([]types.Resource, 0)
	listeners := make([]types.Resource, 0)
	endpoints := make([]types.Resource, 0)
	routes := make([]types.Resource, 0)

	for _, inform := range cp.serviceInformers {
		for _, svc := range inform.GetStore().List() {
			if reflect.TypeOf(svc).Elem().Name() == "Endpoints" {
				continue
			}
			k8sService, ok := svc.(*corev1.Service)
			if !ok {
				cp.log.Errorf("service type is not match, type is: %v", reflect.TypeOf(svc).Elem().Name())
				continue
			}
			// cp.log.Info("=============")
			seviceConfig := ServiceConfig{}
			// seviceConfig.GRPCServiceName = k8sService.Name
			seviceConfig.ServiceName = k8sService.Name
			seviceConfig.Namespace = k8sService.Namespace
			for _, port := range k8sService.Spec.Ports {
				if strings.Contains(port.Name, "grpc") {
					seviceConfig.PortName = port.Name
					seviceConfig.Protocol = "tcp"
					seviceConfig.Region = "us-central1"
					seviceConfig.Zone = "us-central1-a"
					edsService, clsService, rdsService, lsnrService, err := cp.makeXDSConfigFromService(seviceConfig)
					if err != nil {
						cp.log.Errorf("couldn't make service, err: %v", err)
					}
					endpoints = append(endpoints, edsService)
					clusters = append(clusters, clsService)
					routes = append(routes, rdsService)
					listeners = append(listeners, lsnrService)
				}

			}

			// cp.log.Infof("k8sService: %v", k8sService)
			// cp.log.Infof("tmp: %v", seviceConfig)
			// cp.log.Info("=============")
		}
	}

	atomic.AddInt32(&cp.version, 1)
	// cp.log.Infof(" creating snapshot Version " + fmt.Sprint(cp.version))

	// cp.log.Infof("   snapshot with Listener %v", listeners)
	// cp.log.Infof("   snapshot with EDS %v", endpoints)
	// cp.log.Infof("   snapshot with CLS %v", clusters)
	// cp.log.Infof("   snapshot with RDS %v", routes)

	snapshot, err := cachev3.NewSnapshot(fmt.Sprint(cp.version), map[resource.Type][]types.Resource{
		resource.EndpointType: endpoints,
		resource.ClusterType:  clusters,
		resource.ListenerType: listeners,
		resource.RouteType:    routes,
	})
	if err != nil {
		cp.log.Printf(">>>>>>>>>>  Error creating snapshot %v", err)
		return
	}
	IDs := cp.snapshotCache.GetStatusKeys()
	cp.log.Infof("snapshotCache IDs: %v\n", IDs)
	for _, id := range IDs {
		status := cp.snapshotCache.GetStatusInfo(id)
		cp.log.Infof("snapshotCache ID: %v, node meta data: %v", id, status.GetNode().GetMetadata().String())
		err = cp.snapshotCache.SetSnapshot(context.Background(), id, snapshot)
		if err != nil {
			cp.log.Errorf("%v", err)
		}
	}
}
