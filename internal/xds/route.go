package xds

import (
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	corev1 "k8s.io/api/core/v1"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	durationpb "github.com/golang/protobuf/ptypes/duration"
)


func createRoutes(service *corev1.Service) []types.Resource {
	// Create the routes based on the service information
	route := &route.RouteConfiguration{
		Name: "local_route",
		VirtualHosts: []*route.VirtualHost{
			{
				Name:    "local_service",
				Domains: []string{"*"},
				Routes: []*route.Route{
					{
						Match: &route.RouteMatch{
							PathSpecifier: &route.RouteMatch_Prefix{
								Prefix: "/",
							},
						},
						Action: &route.Route_Route{
							Route: &route.RouteAction{
								ClusterSpecifier: &route.RouteAction_Cluster{
									Cluster: "service_cluster",
								},
								Timeout: &durationpb.Duration{
									Seconds: 0,
								},
							},
						},
					},
				},
			},
		},
	}

	return []types.Resource{route}
}