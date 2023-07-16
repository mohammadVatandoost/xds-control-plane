package meshhttproute

import (
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/core"
	api_v1alpha1 "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshhttproute/api/v1alpha1"
	k8s_v1alpha1 "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshhttproute/k8s/v1alpha1"
	plugin_v1alpha1 "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshhttproute/plugin/v1alpha1"
)

func init() {
	core.Register(
		api_v1alpha1.MeshHTTPRouteResourceTypeDescriptor,
		k8s_v1alpha1.AddToScheme,
		plugin_v1alpha1.NewPlugin(),
	)
}
