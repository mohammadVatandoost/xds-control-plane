package informer

import (
	"log/slog"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type ServiceInformer struct {
	cache   cache.SharedIndexInformer
	handler ServiceEventHandler
}

type ServiceEventHandler interface {
	OnAddSerivce(key string, serviceObj *v1.Service)
	OnDeleteService(key string, serviceObj *v1.Service)
	OnUpdateService(newKey string, newServiceObj *v1.Service, oldKey string, oldServiceObj *v1.Service)
}

func NewServiceInformer(factory informers.SharedInformerFactory, handler ServiceEventHandler) *ServiceInformer {
	sharedCache := factory.Core().V1().Services().Informer()

	si := &ServiceInformer{
		cache:   sharedCache,
		handler: handler,
	}

	sharedCache.AddEventHandler(si)

	return si
}

func getServiceKey(service *v1.Service) string {
	return service.Name + "." + service.Namespace
}

func (si *ServiceInformer) Run(stopCh <-chan struct{}) {
	si.cache.Run(stopCh)
}

func (si *ServiceInformer) OnAdd(obj interface{}) {
	service, ok := obj.(*v1.Service)
	if !ok {
		slog.Error("type of object is not service ", "obj", obj, "method", "OnAdd")
		return
	}

	key := getServiceKey(service)
	si.handler.OnAddSerivce(key, service)
}

func (si *ServiceInformer) OnUpdate(oldObj, newObj interface{}) {
	newService, ok := newObj.(*v1.Service)
	if !ok {
		slog.Error("type of object is not service ", "obj", newObj, "method", "OnUpdate")
		return
	}
	oldService, ok := oldObj.(*v1.Service)
	if !ok {
		slog.Error("type of object is not service ", "obj", oldObj, "method", "OnUpdate")
		return
	}
	newKey := getServiceKey(newService)
	oldKey := getServiceKey(oldService)
	si.handler.OnUpdateService(newKey, newService, oldKey, oldService)
}

func (si *ServiceInformer) OnDelete(obj interface{}) {
	service, ok := obj.(*v1.Service)
	if !ok {
		slog.Error("type of object is not service ", "obj", obj, "method", "OnDelete")
		return
	}

	key := getServiceKey(service)
	si.handler.OnDeleteService(key, service)
}
