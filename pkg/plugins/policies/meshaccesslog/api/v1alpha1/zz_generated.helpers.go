// Generated by tools/resource-gen.
// Run "make generate" to update this file.

// nolint:whitespace
package v1alpha1

import (
	common_api "github.com/mohammadVatandoost/xds-conrol-plane/api/common/v1alpha1"
	core_model "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/model"
)

func (x *MeshAccessLog) GetTargetRef() common_api.TargetRef {
	return x.TargetRef
}

func (x *From) GetTargetRef() common_api.TargetRef {
	return x.TargetRef
}

func (x *From) GetDefault() interface{} {
	return x.Default
}

func (x *MeshAccessLog) GetFromList() []core_model.PolicyItem {
	var result []core_model.PolicyItem
	for i := range x.From {
		item := x.From[i]
		result = append(result, &item)
	}
	return result
}

func (x *To) GetTargetRef() common_api.TargetRef {
	return x.TargetRef
}

func (x *To) GetDefault() interface{} {
	return x.Default
}

func (x *MeshAccessLog) GetToList() []core_model.PolicyItem {
	var result []core_model.PolicyItem
	for i := range x.To {
		item := x.To[i]
		result = append(result, &item)
	}
	return result
}
