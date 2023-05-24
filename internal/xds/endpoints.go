package xds

import (
	"context"
	"fmt"
	"reflect"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpointv3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	corev1 "k8s.io/api/core/v1"
	
)

type EnvoyCluster struct {
	name      string
	port      uint32
	endpoints []string
}

func (cp *ControlPlane) HandleServicesUpdate(oldObj, newObj interface{}) {
	cp.log.Info("ControlPlane HandleServicesUpdate")
	for _, inform := range cp.serviceInformers {
		for _, svc := range inform.GetStore().List() {
			
			k8sService, ok := svc.(*corev1.Service)
			if !ok {
				cp.log.Infof("service type is not match, type is: %v", reflect.TypeOf(svc).Elem())
				continue
			}
			cp.log.Infof("Service name: %s, ports: %v\n", k8sService.Name, k8sService.Spec)
		}
	}    
}

func (cp *ControlPlane) HandleEndpointsUpdate(oldObj, newObj interface{}) {
	// cp.log.Info("ControlPlane HandleEndpointsUpdate")

	edsServiceData := make(map[string]*EnvoyCluster, 0)
	// rt := make([]types.Resource, 0)
	// sec := make([]types.Resource, 0)
	// lbEndPoints := make([]types.Resource, 0)
	// lbEndPoints := make([]types.Resource, 0)
	// lbEndPoints := make([]types.Resource, 0)

	for _, inform := range cp.endpointInformers {
		for _, ep := range inform.GetStore().List() {

			endpoints := ep.(*corev1.Endpoints)
			// cp.log.Infof("endpoints Labels: %v", endpoints.Labels)
			//ToDo: use it only for specefic services
			if _, ok := endpoints.Labels["xds"]; !ok {
				continue
			}

			if _, ok := edsServiceData[endpoints.Name]; !ok {
				edsServiceData[endpoints.Name] = &EnvoyCluster{
					name: endpoints.Name,
				}
			}
			// cp.log.Infof("endpoints: %v", endpoints.String())
			for _, subset := range endpoints.Subsets {
				cp.log.Infof("endpoints subset: %v", subset.String())
				for i, addr := range subset.Addresses {
					// cp.log.Infof("endpoints Subsets addresses, IP: %v, Port: %v", addr.IP, subset.Ports[i].Port)
					edsServiceData[endpoints.Name].port = uint32(subset.Ports[i].Port)
					edsServiceData[endpoints.Name].endpoints = append(edsServiceData[endpoints.Name].endpoints, addr.IP)
				}
			}
		}
	}

	// for each pod create endpoints
	edsEndpoints := make([]types.Resource, len(edsServiceData))
	for _, envoyCluster := range edsServiceData {
		edsEndpoints = append(edsEndpoints, cp.MakeEndpointsForCluster(envoyCluster))
	}

	//snapshot := cache.NewSnapshot(fmt.Sprintf("%v.0", version), edsEndpoints, nil, nil, nil, nil, nil)
	snapshot, err := cache.NewSnapshot(fmt.Sprintf("%v.0", cp.version), map[resource.Type][]types.Resource{
		resource.EndpointType: edsEndpoints,
	})
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	IDs := cp.snapshotCache.GetStatusKeys()
	for _, id := range IDs {
		err = cp.snapshotCache.SetSnapshot(context.Background(), id, snapshot)
		if err != nil {
			fmt.Printf("%v", err)
		}
	}

	cp.version++
}

func (cp *ControlPlane) MakeEndpointsForCluster(service *EnvoyCluster) *endpointv3.ClusterLoadAssignment {
	// cp.log.Infof("Updating endpoints for cluster,  service.name %s: service.endpoints: %v", service.name, service.endpoints)
	cla := &endpointv3.ClusterLoadAssignment{
		ClusterName: service.name,
		Endpoints:   []*endpointv3.LocalityLbEndpoints{},
	}

	for _, endpoint := range service.endpoints {
		cla.Endpoints = append(cla.Endpoints,
			&endpointv3.LocalityLbEndpoints{
				LbEndpoints: []*endpointv3.LbEndpoint{{
					HostIdentifier: &endpointv3.LbEndpoint_Endpoint{
						Endpoint: &endpointv3.Endpoint{
							Address: &core.Address{
								Address: &core.Address_SocketAddress{
									SocketAddress: &core.SocketAddress{
										Protocol: core.SocketAddress_TCP,
										Address:  endpoint,
										PortSpecifier: &core.SocketAddress_PortValue{
											PortValue: service.port,
										},
									},
								},
							},
						},
					},
				}},
			},
		)
	}
	return cla
}
