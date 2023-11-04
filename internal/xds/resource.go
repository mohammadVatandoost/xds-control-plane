package xds

import (
	"fmt"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	v3routerpb "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/mohammadVatandoost/xds-conrol-plane/internal/resource"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func MakeXDSResource(resourceInfo *resource.Resource, region string,
	zone string, portName string) (*endpoint.ClusterLoadAssignment, *cluster.Cluster, *route.RouteConfiguration, *listener.Listener, error) {
	routeConfigName := resourceInfo.Name + "-route"
	clusterName := resourceInfo.Name + "-cluster"
	virtualHostName := resourceInfo.Name + "-vs"
	addresses := getAddresses(resourceInfo.Key, portName)
	if len(addresses) == 0 {
		return nil, nil, nil, nil, fmt.Errorf("there is no availabe address for service: %v", resourceInfo.Key)
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
		Domains: []string{resourceInfo.Key}, //******************* >> must match what is specified at xds:/// //

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

	filterPbst, err := anypb.New(&v3routerpb.Router{})
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to marshal the router, key: %v, err: %v",
			resourceInfo.Key, err)
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

	pbst, err := anypb.New(manager)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to marshal the manager message, key: %v, err: %v",
			resourceInfo.Key, err)
	}

	lsnr := &listener.Listener{
		Name: resourceInfo.Key,
		ApiListener: &listener.ApiListener{
			ApiListener: pbst,
		},
	}

	return eds, cls, rds, lsnr, nil
}
