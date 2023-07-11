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

	// clusters := make([]types.Resource, 0)
	// listeners := make([]types.Resource, 0)
	// endpoints := make([]types.Resource, 0)
	// routes := make([]types.Resource, 0)

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
					nodes := cp.GetNodesWatchTheResource(seviceConfig.ServiceName)
					cp.log.Infof("ControlPlane HandleServicesUpdate nodes: %v, serviceName: %s", nodes, seviceConfig.ServiceName)
					for _, n := range nodes {
						node, err := cp.GetNode(n)
						if err != nil {
							cp.log.Errorf("node id: %s is not watched the resource id: %s", n, seviceConfig.ServiceName)
							continue
						}
						cp.log.Infof("ControlPlane HandleServicesUpdate nodes: %v, serviceName: %s, listeners: %v", nodes, seviceConfig.ServiceName, lsnrService)
						node.AddCluster(clsService)
						node.AddEndpoint(edsService)
						node.AddRoute(rdsService)
						node.AddListener(lsnrService)
					}
				}
			}
		}
	}

	atomic.AddInt32(&cp.version, 1)

	IDs := cp.snapshotCache.GetStatusKeys()
	cp.log.Infof("snapshotCache IDs: %v\n", IDs)
	for _, id := range IDs {
		node, err := cp.GetNode(id)
		if err != nil {
			cp.log.Errorf("node id: %s is not exist", id)
			continue
		}
		snapshot, err := cachev3.NewSnapshot(fmt.Sprint(cp.version), map[resource.Type][]types.Resource{
			resource.EndpointType: node.GetEndpoints(),
			resource.ClusterType:  node.GetClusters(),
			resource.ListenerType: node.GetListeners(),
			resource.RouteType:    node.GetRoutes(),
		})
		if err != nil {
			cp.log.Printf(">>>>>>>>>>  Error creating snapshot %v", err)
			return
		}
		status := cp.snapshotCache.GetStatusInfo(id)
		cp.log.Infof("snapshotCache ID: %v, node meta data: %v", id, status.GetNode().GetMetadata().String())
		err = cp.snapshotCache.SetSnapshot(context.Background(), id, snapshot)
		if err != nil {
			cp.log.Errorf("%v", err)
		}
		node.ClearResources()
	}
}
