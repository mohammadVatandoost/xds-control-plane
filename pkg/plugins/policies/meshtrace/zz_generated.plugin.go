package meshtrace

import (
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/core"
	api_v1alpha1 "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshtrace/api/v1alpha1"
	k8s_v1alpha1 "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshtrace/k8s/v1alpha1"
	plugin_v1alpha1 "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshtrace/plugin/v1alpha1"
)

func init() {
	core.Register(
		api_v1alpha1.MeshTraceResourceTypeDescriptor,
		k8s_v1alpha1.AddToScheme,
		plugin_v1alpha1.NewPlugin(),
	)
}
