package xds

import (
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/golang/protobuf/ptypes"
	corev1 "k8s.io/api/core/v1"
)

func createClusters(service *corev1.Service) []types.Resource {
	// Create the clusters based on the service information
	cluster := &cluster.Cluster{
		Name:           "service_cluster",
		ConnectTimeout: ptypes.DurationProto(5 * time.Second),
		ClusterDiscoveryType: &cluster.Cluster_Type{
			Type: cluster.Cluster_STRICT_DNS,
		},
		LbPolicy: cluster.Cluster_ROUND_ROBIN,
		LoadAssignment: &endpoint.ClusterLoadAssignment{
			ClusterName: "service_cluster",
			Endpoints: []*endpoint.LocalityLbEndpoints{
				{
					LbEndpoints: []*endpoint.LbEndpoint{
						{
							HostIdentifier: &endpoint.LbEndpoint_Endpoint{
								Endpoint: &endpoint.Endpoint{
									Address: &core.Address{
										Address: &core.Address_SocketAddress{
											SocketAddress: &core.SocketAddress{
												Protocol: core.SocketAddress_TCP,
												Address:  service.Name + "." + service.Namespace + ".svc.cluster.local",
												PortSpecifier: &core.SocketAddress_PortValue{
													PortValue: uint32(service.Spec.Ports[0].Port),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return []types.Resource{cluster}
}
