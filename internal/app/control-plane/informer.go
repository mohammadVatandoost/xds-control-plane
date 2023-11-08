package controlplane

import (
	"log/slog"

	"github.com/mohammadVatandoost/xds-conrol-plane/internal/resource"
)

func (a *App) OnAddSerivce(res *resource.Resource) {
	a.muResource.Lock()
	defer a.muResource.Unlock()
	slog.Info("OnAddSerivce", "key", res.Key, "name", res.ServiceObj.Name,
		"Namespace", res.ServiceObj.Namespace, "Labels", res.ServiceObj.Labels)
	a.resources[res.Key] = res
}

// func (a *App) OnDeleteService(key string, serviceObj *v1.Service) {
// 	a.muResource.Lock()
// 	defer a.muResource.Unlock()
// 	slog.Info("OnDeleteService", "key", key, "name", serviceObj.Name, "Namespace", serviceObj.Namespace, "Labels", serviceObj.Labels)

// 	_, ok := a.resources[key]
// 	if !ok {
// 		slog.Error("OnDeleteService resource doesn't exist in DB", "key", key, "name", serviceObj.Name, "Namespace", serviceObj.Namespace, "Labels", serviceObj.Labels)
// 		return
// 	}
// 	delete(a.resources, key)
// }

func (a *App) DeleteService(key string) {
	a.muResource.Lock()
	defer a.muResource.Unlock()
	delete(a.resources, key)
}

func (a *App) OnUpdateService(newRes *resource.Resource, oldRes *resource.Resource) {
	// slog.Info("OnUpdateService", "newKey", newKey, "newServiceName", newServiceObj.Name, "newServiceNamespace", newServiceObj.Namespace,
	// 	"oldKey", oldKey, "oleServiceName", oldServiceObj.Name, "oldServiceNamespace", oldServiceObj.Namespace)
	a.muResource.Lock()
	defer a.muResource.Unlock()
	_, ok := a.resources[oldRes.Key]
	if !ok {
		slog.Error("OnUpdateService resource doesn't exist in DB", "key", oldRes.Key, "name", oldRes.ServiceObj.Name,
			"Namespace", oldRes.ServiceObj.Namespace, "Labels", oldRes.ServiceObj.Labels)
	}
	delete(a.resources, oldRes.Key)
	a.resources[newRes.Key] = newRes
}
