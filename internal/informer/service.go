package informer

import (
	"log/slog"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type ServiceInformer struct {
	cache cache.SharedIndexInformer
}

type ServiceEventHandler interface {
}

func NewServiceInformer(factory informers.SharedInformerFactory) *ServiceInformer {
	sharedCache := factory.Core().V1().Services().Informer()

	si := &ServiceInformer{
		cache: sharedCache,
	}

	sharedCache.AddEventHandler(si)

	return si
}

func getServiceKey(service *v1.Service) string {
	return service.Namespace + "." + service.Name
}

func (si *ServiceInformer) OnAdd(obj interface{}) {
	service, ok := obj.(*v1.Service)
	if !ok {
		slog.Error("type of object is not service ", "obj", obj, "method", "OnAdd")
		return
	}

	key := getServiceKey(service)
	si.OnAddSerivce(key, service)
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
	si.OnUpdateService(newKey, newService, oldKey, oldService)
}

func (si *ServiceInformer) OnDelete(obj interface{}) {
	service, ok := obj.(*v1.Service)
	if !ok {
		slog.Error("type of object is not service ", "obj", obj, "method", "OnDelete")
		return
	}

	key := getServiceKey(service)
	si.OnDeleteService(key, service)
}
