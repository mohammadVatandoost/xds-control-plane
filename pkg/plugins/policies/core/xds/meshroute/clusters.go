package meshroute

import (
	"github.com/pkg/errors"

	mesh_proto "github.com/mohammadVatandoost/xds-conrol-plane/api/mesh/v1alpha1"
	core_mesh "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/apis/mesh"
	core_xds "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/xds"
	xds_context "github.com/mohammadVatandoost/xds-conrol-plane/pkg/xds/context"
	envoy_common "github.com/mohammadVatandoost/xds-conrol-plane/pkg/xds/envoy"
	envoy_clusters "github.com/mohammadVatandoost/xds-conrol-plane/pkg/xds/envoy/clusters"
	envoy_tags "github.com/mohammadVatandoost/xds-conrol-plane/pkg/xds/envoy/tags"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/xds/generator"
)

func GenerateClusters(
	proxy *core_xds.Proxy,
	meshCtx xds_context.MeshContext,
	services envoy_common.Services,
) (*core_xds.ResourceSet, error) {
	resources := core_xds.NewResourceSet()

	for _, serviceName := range services.Sorted() {
		service := services[serviceName]
		protocol := generator.InferProtocol(proxy, service.Clusters())
		tlsReady := service.TLSReady()

		for _, cluster := range service.Clusters() {
			edsClusterBuilder := envoy_clusters.NewClusterBuilder(proxy.APIVersion)

			clusterName := cluster.Name()
			clusterTags := []envoy_tags.Tags{cluster.Tags()}

			if service.HasExternalService() {
				if meshCtx.Resource.ZoneEgressEnabled() {
					edsClusterBuilder.
						Configure(envoy_clusters.EdsCluster(clusterName)).
						Configure(envoy_clusters.ClientSideMTLS(
							proxy.SecretsTracker,
							meshCtx.Resource,
							mesh_proto.ZoneEgressServiceName,
							tlsReady,
							clusterTags,
						))
				} else {
					endpoints := proxy.Routing.ExternalServiceOutboundTargets[serviceName]
					isIPv6 := proxy.Dataplane.IsIPv6()

					edsClusterBuilder.
						Configure(envoy_clusters.ProvidedEndpointCluster(clusterName, isIPv6, endpoints...)).
						Configure(envoy_clusters.ClientSideTLS(endpoints))
				}

				switch protocol {
				case core_mesh.ProtocolHTTP:
					edsClusterBuilder.Configure(envoy_clusters.Http())
				case core_mesh.ProtocolHTTP2, core_mesh.ProtocolGRPC:
					edsClusterBuilder.Configure(envoy_clusters.Http2())
				default:
				}
			} else {
				edsClusterBuilder.
					Configure(envoy_clusters.EdsCluster(clusterName)).
					Configure(envoy_clusters.Http2())

				if upstreamMeshName := cluster.Mesh(); upstreamMeshName != "" {
					for _, otherMesh := range append(meshCtx.Resources.OtherMeshes().Items, meshCtx.Resource) {
						if otherMesh.GetMeta().GetName() == upstreamMeshName {
							edsClusterBuilder.Configure(
								envoy_clusters.CrossMeshClientSideMTLS(
									proxy.SecretsTracker, meshCtx.Resource, otherMesh, serviceName, tlsReady, clusterTags,
								),
							)
							break
						}
					}
				} else {
					edsClusterBuilder.Configure(envoy_clusters.ClientSideMTLS(
						proxy.SecretsTracker,
						meshCtx.Resource, serviceName, tlsReady, clusterTags))
				}
			}

			edsCluster, err := edsClusterBuilder.Build()
			if err != nil {
				return nil, errors.Wrapf(err, "build CDS for cluster %s failed", clusterName)
			}

			resources = resources.Add(&core_xds.Resource{
				Name:     clusterName,
				Origin:   generator.OriginOutbound,
				Resource: edsCluster,
			})
		}
	}

	return resources, nil
}
