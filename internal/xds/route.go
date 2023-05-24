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
		Name: service.Name + "." + service.Namespace + ".svc.cluster.local",
		VirtualHosts: []*route.VirtualHost{
			{
				Name:    service.Name + "." + service.Namespace + ".svc.cluster.local",
				Domains: []string{service.Name + "." + service.Namespace + ".svc.cluster.local", service.Name, service.Name + "." + service.Namespace},
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
									Cluster: service.Name+"-cluster",
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