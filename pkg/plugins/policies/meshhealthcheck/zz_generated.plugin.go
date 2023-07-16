package meshhealthcheck

import (
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/core"
	api_v1alpha1 "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshhealthcheck/api/v1alpha1"
	k8s_v1alpha1 "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshhealthcheck/k8s/v1alpha1"
	plugin_v1alpha1 "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshhealthcheck/plugin/v1alpha1"
)

func init() {
	core.Register(
		api_v1alpha1.MeshHealthCheckResourceTypeDescriptor,
		k8s_v1alpha1.AddToScheme,
		plugin_v1alpha1.NewPlugin(),
	)
}
