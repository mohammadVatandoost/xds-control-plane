package rest

import (
	core_model "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/model"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/model/rest/v1alpha1"
)

type Resource interface {
	GetMeta() v1alpha1.ResourceMeta
	GetSpec() core_model.ResourceSpec
}
