package xds

import (
	"fmt"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/mohammadVatandoost/xds-conrol-plane/internal/resource"
)

func (cp *ControlPlane) makeXDSResource(resourceInfo *resource.Resource, region string, zone string) (*endpoint.ClusterLoadAssignment, *cluster.Cluster, *route.RouteConfiguration, *listener.Listener, error) {
	routeConfigName := resourceInfo.Name + "-route"
	clusterName := resourceInfo.Name + "-cluster"
	virtualHostName := resourceInfo.Name + "-vs"
	addresses := getAddresses(resourceInfo.Key)
	if len(addresses) == 0 {
		return nil, nil, nil, nil, fmt.Errorf("there is no availabe address for service: %v", svc.ServiceName)
	}
	// cp.log.Infof("service: %v, addresses: %v \n", svc.ServiceName, addresses)
	lbe := makeLBEndpoint(addresses)
	eds := &endpoint.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpoint.LocalityLbEndpoints{{
			Locality: &core.Locality{
				Region: region,
				Zone:   zone,
			},
			Priority:            0,
			LoadBalancingWeight: &wrapperspb.UInt32Value{Value: uint32(1000)},
			LbEndpoints:         lbe,
		}},
	}
	cls := &cluster.Cluster{
		Name:                 clusterName,
		LbPolicy:             cluster.Cluster_ROUND_ROBIN,
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_EDS},
		EdsClusterConfig: &cluster.Cluster_EdsClusterConfig{
			EdsConfig: &core.ConfigSource{
				ConfigSourceSpecifier: &core.ConfigSource_Ads{},
			},
		},
	}

	// RDS
	// cp.log.Infof(">>>>>>>>>>>>>>>>>>> creating RDS " + virtualHostName)
	vh := &route.VirtualHost{
		Name:    virtualHostName,
		Domains: []string{svc.ServiceName}, //******************* >> must match what is specified at xds:/// //

		Routes: []*route.Route{{
			Match: &route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Prefix{
					Prefix: "",
				},
			},
			Action: &route.Route_Route{
				Route: &route.RouteAction{
					ClusterSpecifier: &route.RouteAction_Cluster{
						Cluster: clusterName,
					},
				},
			},
		}}}

	rds := &route.RouteConfiguration{
		Name:         routeConfigName,
		VirtualHosts: []*route.VirtualHost{vh},
	}

	// LISTENER
	// cp.log.Infof(">>>>>>>>>>>>>>>>>>> creating LISTENER " + svc.ServiceName)
	hcRds := &hcm.HttpConnectionManager_Rds{
		Rds: &hcm.Rds{
			RouteConfigName: routeConfigName,
			ConfigSource: &core.ConfigSource{
				ConfigSourceSpecifier: &core.ConfigSource_Ads{
					Ads: &core.AggregatedConfigSource{},
				},
			},
		},
	}

	filterPbst, err := ptypes.MarshalAny(&v3routerpb.Router{})
	if err != nil {
		panic(err)
	}
	// RouterHTTPFilter := hcm.HTTPFilter("router", &v3routerpb.Router{})
	RouterHTTPFilter := &hcm.HttpFilter{
		Name: "router",
		ConfigType: &hcm.HttpFilter_TypedConfig{
			TypedConfig: filterPbst,
		},
	}
	filters := []*hcm.HttpFilter{
		RouterHTTPFilter,
	}

	manager := &hcm.HttpConnectionManager{
		CodecType:      hcm.HttpConnectionManager_AUTO,
		RouteSpecifier: hcRds,
		HttpFilters:    filters,
	}

	pbst, err := ptypes.MarshalAny(manager)
	if err != nil {
		panic(err)
	}

	lsnr := &listener.Listener{
		Name: svc.ServiceName,
		ApiListener: &listener.ApiListener{
			ApiListener: pbst,
		},
	}

	return eds, cls, rds, lsnr, nil
}