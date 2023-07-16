package meshratelimit

import (
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/core"
	api_v1alpha1 "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshratelimit/api/v1alpha1"
	k8s_v1alpha1 "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshratelimit/k8s/v1alpha1"
	plugin_v1alpha1 "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshratelimit/plugin/v1alpha1"
)

func init() {
	core.Register(
		api_v1alpha1.MeshRateLimitResourceTypeDescriptor,
		k8s_v1alpha1.AddToScheme,
		plugin_v1alpha1.NewPlugin(),
	)
}
