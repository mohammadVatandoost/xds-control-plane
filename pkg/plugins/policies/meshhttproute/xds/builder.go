package xds

import (
	envoy_listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"

	common_api "github.com/mohammadVatandoost/xds-conrol-plane/api/common/v1alpha1"
	mesh_proto "github.com/mohammadVatandoost/xds-conrol-plane/api/mesh/v1alpha1"
	api "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshhttproute/api/v1alpha1"
	envoy_common "github.com/mohammadVatandoost/xds-conrol-plane/pkg/xds/envoy"
	envoy_listeners_v3 "github.com/mohammadVatandoost/xds-conrol-plane/pkg/xds/envoy/listeners/v3"
	envoy_names "github.com/mohammadVatandoost/xds-conrol-plane/pkg/xds/envoy/names"
	envoy_routes "github.com/mohammadVatandoost/xds-conrol-plane/pkg/xds/envoy/routes"
	envoy_virtual_hosts "github.com/mohammadVatandoost/xds-conrol-plane/pkg/xds/envoy/virtualhosts"
)

type OutboundRoute struct {
	Matches                 []api.Match
	Filters                 []api.Filter
	Split                   []envoy_common.Split
	BackendRefToClusterName map[common_api.TargetRefHash]string
}

type HttpOutboundRouteConfigurer struct {
	Service string
	Routes  []OutboundRoute
	DpTags  mesh_proto.MultiValueTagSet
}

var _ envoy_listeners_v3.FilterChainConfigurer = &HttpOutboundRouteConfigurer{}

func (c *HttpOutboundRouteConfigurer) Configure(filterChain *envoy_listener.FilterChain) error {
	virtualHostBuilder := envoy_virtual_hosts.NewVirtualHostBuilder(envoy_common.APIV3).
		Configure(envoy_virtual_hosts.CommonVirtualHost(c.Service))
	for _, route := range c.Routes {
		route := envoy_virtual_hosts.AddVirtualHostConfigurer(
			&RoutesConfigurer{
				Matches:                 route.Matches,
				Filters:                 route.Filters,
				Split:                   route.Split,
				BackendRefToClusterName: route.BackendRefToClusterName,
			})
		virtualHostBuilder = virtualHostBuilder.Configure(route)
	}
	static := envoy_listeners_v3.HttpStaticRouteConfigurer{
		Builder: envoy_routes.NewRouteConfigurationBuilder(envoy_common.APIV3).
			Configure(envoy_routes.CommonRouteConfiguration(envoy_names.GetOutboundRouteName(c.Service))).
			Configure(envoy_routes.TagsHeader(c.DpTags)).
			Configure(envoy_routes.VirtualHost(virtualHostBuilder)),
	}

	return static.Configure(filterChain)
}
