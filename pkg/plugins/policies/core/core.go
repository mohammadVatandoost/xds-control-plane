package core

import (
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/plugins"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/model"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/registry"
)

var AllSchemes []func(*runtime.Scheme) error

func Register(resType model.ResourceTypeDescriptor, fn func(scheme *runtime.Scheme) error, p plugins.Plugin) {
	plugins.Register(plugins.PluginName(resType.KumactlArg), p)
	registry.RegisterType(resType)
	AllSchemes = append(AllSchemes, fn)
}
