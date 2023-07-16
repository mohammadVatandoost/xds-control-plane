package v1alpha1

import (
	"github.com/pkg/errors"

	core_plugins "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/plugins"
	core_mesh "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/apis/mesh"
	core_xds "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/xds"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/core/matchers"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/core/xds/meshroute"
	api "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshtcproute/api/v1alpha1"
	xds_context "github.com/mohammadVatandoost/xds-conrol-plane/pkg/xds/context"
	envoy_common "github.com/mohammadVatandoost/xds-conrol-plane/pkg/xds/envoy"
)

var _ core_plugins.PolicyPlugin = &plugin{}

type plugin struct{}

func NewPlugin() core_plugins.Plugin {
	return &plugin{}
}

func (p plugin) MatchedPolicies(
	dataplane *core_mesh.DataplaneResource,
	resources xds_context.Resources,
) (core_xds.TypedMatchingPolicies, error) {
	return matchers.MatchedPolicies(api.MeshTCPRouteType, dataplane, resources)
}

func (p plugin) Apply(
	rs *core_xds.ResourceSet,
	ctx xds_context.Context,
	proxy *core_xds.Proxy,
) error {
	tcpRules := proxy.Policies.Dynamic[api.MeshTCPRouteType].ToRules.Rules
	if len(tcpRules) == 0 {
		return nil
	}

	servicesAccumulator := envoy_common.NewServicesAccumulator(
		ctx.Mesh.ServiceTLSReadiness)

	listeners, err := generateListeners(proxy, tcpRules, servicesAccumulator)
	if err != nil {
		return errors.Wrap(err, "couldn't generate listener resources")
	}
	rs.AddSet(listeners)

	services := servicesAccumulator.Services()

	clusters, err := meshroute.GenerateClusters(proxy, ctx.Mesh, services)
	if err != nil {
		return errors.Wrap(err, "couldn't generate cluster resources")
	}
	rs.AddSet(clusters)

	endpoints, err := meshroute.GenerateEndpoints(proxy, ctx, services)
	if err != nil {
		return errors.Wrap(err, "couldn't generate endpoint resources")
	}
	rs.AddSet(endpoints)

	return nil
}
