package resource

import v1 "k8s.io/api/core/v1"

type Resource struct {
	Name         string              `json:"name"`
	Version      string              `json:"version"`
	NameSpace    string              `json:"namespace"`
	K8SKind      string              `json:"kind"`
	EnvoyTypeURL string              `json:"envoyTypeURL"`
	Key          string              `json:"key"` //key is name.namespace:portnumber
	Watchers     map[string]struct{} `json:"watchers"`
	ServiceObj   *v1.Service         `json:"serviceOBJ"` // we only support Service, later we can add Ingress
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
