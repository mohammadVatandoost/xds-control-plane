package xds

import (
	"fmt"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/golang/protobuf/ptypes/any"
	corev1 "k8s.io/api/core/v1"
)

func createListeners(service *corev1.Service) []types.Resource {
	// Create the listeners based on the service information
	// serviceName := service.Name + "." + service.Namespace + ".svc.cluster.local"
	serviceName := fmt.Sprintf("%s:%d", service.Name, service.Spec.Ports[0].Port)
	listener := &listener.Listener{
		Name: serviceName,
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: uint32(service.Spec.Ports[0].Port),
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{
			{
				Filters: []*listener.Filter{
					{
						Name: "envoy.filters.network.http_connection_manager",
						ConfigType: &listener.Filter_TypedConfig{
							TypedConfig: &any.Any{
								Value: []byte(`{
									"stat_prefix": "ingress_http",
									"route_config": {
										"name": "local_route",
										"virtual_hosts": [
											{
												"name": "local_service",
												"domains": ["*"],
												"routes": [
													{
														"match": {
															"prefix": "/"
														},
														"route": {
															"cluster": "service_cluster",
															"timeout": "0s"
														}
													}
												]
											}
										]
									},
									"http_filters": [
										{
											"name": "envoy.filters.http.router"
										}
									]
								}`),
							},
						},
					},
				},
			},
		},
	}

	return []types.Resource{listener}
}
