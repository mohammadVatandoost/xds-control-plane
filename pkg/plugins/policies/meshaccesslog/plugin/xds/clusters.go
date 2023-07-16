package xds

import (
	core_xds "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/xds"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/xds/envoy/clusters"
)

const OriginMeshAccessLog = "meshaccesslog"

type LoggingEndpoint struct {
	Address string
	Port    uint32
}

func xdsEndpoint(endpoint LoggingEndpoint) core_xds.Endpoint {
	return core_xds.Endpoint{
		Target: endpoint.Address,
		Port:   endpoint.Port,
	}
}

func HandleClusters(backendEndpoints EndpointAccumulator, rs *core_xds.ResourceSet, proxy *core_xds.Proxy) error {
	for backendEndpoint := range backendEndpoints.endpoints {
		endpoint := xdsEndpoint(backendEndpoint)

		clusterName := backendEndpoints.clusterForEndpoint(backendEndpoint)
		res, err := clusters.NewClusterBuilder(proxy.APIVersion).
			Configure(clusters.ProvidedEndpointCluster(string(clusterName), proxy.Dataplane.IsIPv6(), endpoint)).
			Configure(clusters.ClientSideTLS([]core_xds.Endpoint{endpoint})).
			Configure(clusters.DefaultTimeout()).
			Configure(clusters.Http2()).
			Build()
		if err != nil {
			return err
		}

		rs.Add(&core_xds.Resource{Name: string(clusterName), Origin: OriginMeshAccessLog, Resource: res})
	}

	return nil
}
