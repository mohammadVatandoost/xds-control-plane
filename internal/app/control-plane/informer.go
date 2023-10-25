package controlplane

import (
	"log/slog"

	"github.com/mohammadVatandoost/xds-conrol-plane/internal/resource"
	v1 "k8s.io/api/core/v1"
)

func (a *App) OnAddSerivce(key string, serviceObj *v1.Service) {
	a.muResource.Lock()
	defer a.muResource.Unlock()
	slog.Info("OnAddSerivce", "key", key, "name", serviceObj.Name, "Namespace", serviceObj.Namespace, "Labels", serviceObj.Labels)
	resourceInstance, ok := a.resources[key]
	if !ok {
		resourceInstance = resource.NewResource(serviceObj.Name, serviceObj.APIVersion, "", "service", key)
	}
	resourceInstance.Name = serviceObj.Name
	resourceInstance.Version = serviceObj.APIVersion
	a.resources[key] = resourceInstance
}

func (a *App) OnDeleteService(key string, serviceObj *v1.Service) {
	a.muResource.Lock()
	defer a.muResource.Unlock()
	slog.Info("OnDeleteService", "key", key, "name", serviceObj.Name, "Namespace", serviceObj.Namespace, "Labels", serviceObj.Labels)

	_, ok := a.resources[key]
	if !ok {
		slog.Error("OnDeleteService resource doesn't exist in DB", "key", key, "name", serviceObj.Name, "Namespace", serviceObj.Namespace, "Labels", serviceObj.Labels)
		return
	}
	delete(a.resources, key)
	
}

func (a *App) OnUpdateService(newKey string, newServiceObj *v1.Service, oldKey string, oldServiceObj *v1.Service) {
	slog.Info("OnUpdateService", "newKey", newKey, "newServiceName", newServiceObj.Name, "newServiceNamespace", newServiceObj.Namespace,
	"oldKey", oldKey, "oleServiceName", oldServiceObj.Name, "oldServiceNamespace", oldServiceObj.Namespace)
	a.muResource.Lock()
	defer a.muResource.Unlock()
	resourceInstance, ok := a.resources[oldKey]
	if !ok {
		slog.Error("OnUpdateService resource doesn't exist in DB", "key", oldKey, "name", oldServiceObj.Name,
		 "Namespace", oldServiceObj.Namespace, "Labels", oldServiceObj.Labels)
		resourceInstance = resource.NewResource(newServiceObj.Name, newServiceObj.APIVersion, "", "service", newKey)
	}
	delete(a.resources, oldKey)
	resourceInstance.Name = newServiceObj.Name
	resourceInstance.Version = newServiceObj.APIVersion
	a.resources[newKey] = resourceInstance
}