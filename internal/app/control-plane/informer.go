package controlplane

import (
	"log/slog"

	v1 "k8s.io/api/core/v1"
)

func (a *App) OnAddSerivce(key string, serviceObj *v1.Service) {
	slog.Info("OnAddSerivce", "key", key, "name", serviceObj.Name, "Namespace", serviceObj.Namespace, "Labels", serviceObj.Labels)
}

func (a *App) OnDeleteService(key string, serviceObj *v1.Service) {
	slog.Info("OnDeleteService", "key", key, "name", serviceObj.Name, "Namespace", serviceObj.Namespace, "Labels", serviceObj.Labels)
}

func (a *App) OnUpdateService(newKey string, newServiceObj *v1.Service, oldKey string, oldServiceObj *v1.Service) {
	slog.Info("OnUpdateService", "newKey", newKey, "newServiceName", newServiceObj.Name, "newServiceNamespace", newServiceObj.Namespace,
	"oldKey", oldKey, "oleServiceName", oldServiceObj.Name, "oldServiceNamespace", oldServiceObj.Namespace)
}