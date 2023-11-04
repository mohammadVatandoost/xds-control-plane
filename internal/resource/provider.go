package resource

import v1 "k8s.io/api/core/v1"

type Resource struct {
	Name         string
	Version      string
	NameSpace    string
	K8SKind      string
	EnvoyTypeURL string
	Key          string //key is name.namespace:portnumber
	Watchers     map[string]struct{}
	ServiceObj   *v1.Service // we only support Service, later we can add Ingress
}

func NewResource(name, version, nameSpace, resourceType, key string, serviceObj *v1.Service) *Resource {
	return &Resource{
		Name:         name,
		Version:      version,
		NameSpace:    nameSpace,
		EnvoyTypeURL: resourceType,
		Key:          key,
		ServiceObj:   serviceObj,
		Watchers:     make(map[string]struct{}),
	}
}
