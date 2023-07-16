package v1alpha1

import (
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/core"
	core_plugins "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/plugins"
	core_mesh "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/apis/mesh"
	core_xds "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/xds"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/core/matchers"
	api "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/donothingpolicy/api/v1alpha1"
	xds_context "github.com/mohammadVatandoost/xds-conrol-plane/pkg/xds/context"
)

var (
	_   core_plugins.PolicyPlugin = &plugin{}
	log                           = core.Log.WithName("DoNothingPolicy")
)

type plugin struct{}

func NewPlugin() core_plugins.Plugin {
	return &plugin{}
}

func (p plugin) MatchedPolicies(dataplane *core_mesh.DataplaneResource, resources xds_context.Resources) (core_xds.TypedMatchingPolicies, error) {
	return matchers.MatchedPolicies(api.DoNothingPolicyType, dataplane, resources)
}

func (p plugin) Apply(rs *core_xds.ResourceSet, ctx xds_context.Context, proxy *core_xds.Proxy) error {
	log.Info("apply is not implemented")
	return nil
}
