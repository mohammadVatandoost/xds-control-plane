package informer

import (
	"fmt"
	"log/slog"

	"github.com/mohammadVatandoost/xds-conrol-plane/internal/resource"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

const PortNameLabel = "xds/portName"

type ServiceInformer struct {
	cache   cache.SharedIndexInformer
	handler ServiceEventHandler
}

type ServiceEventHandler interface {
	OnAddSerivce(res *resource.Resource)
	OnUpdateService(newRes *resource.Resource, oldRes *resource.Resource)
	DeleteService(key string)
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

// ToDo: later use array of ports
func isXDSService(service *v1.Service) (string, bool) {
	for k, value := range service.Labels {
		if k == PortNameLabel {
			return value, true
		}
	}
	return "", false
}

func getServiceKey(service *v1.Service, portName string) (string, error) {
	for _, port := range service.Spec.Ports {
		if port.Name == portName {
			return fmt.Sprintf("%s.%s.svc.cluster.local:%d",
				service.Name, service.Namespace, port.Port), nil
		}
	}
	return "", fmt.Errorf("couldn't find the port name, portName: %s, serviceName: %s, namespace: %s",
		portName, service.Name, service.Namespace)
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
	portName, ok := isXDSService(service)
	if !ok {
		return
	}
	key, err := getServiceKey(service, portName)
	if err != nil {
		slog.Error("couldn't get service key ", "err", err)
		return
	}
	res := resource.NewResource(service.Name, service.APIVersion, "", "service", key, portName, service)
	si.handler.OnAddSerivce(res)
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

	portNameOld, ok := isXDSService(oldService)
	if !ok {
		portNameNew, ok := isXDSService(newService)
		if ok {
			key, err := getServiceKey(newService, portNameNew)
			if err != nil {
				slog.Error("couldn't get service key ", "err", err)
				return
			}
			res := resource.NewResource(newService.Name, newService.APIVersion, "", "service", key, portNameNew, newService)
			si.handler.OnAddSerivce(res)
		}
		return
	}

	portNameNew, ok := isXDSService(newService)
	if !ok {
		oldKey, err := getServiceKey(oldService, portNameOld)
		if err != nil {
			slog.Error("couldn't get service key ", "err", err)
			return
		}
		si.handler.DeleteService(oldKey)
		return
	}

	oldKey, err := getServiceKey(oldService, portNameOld)
	if err != nil {
		slog.Error("couldn't get oldService key ", "err", err)
		return
	}

	newKey, err := getServiceKey(newService, portNameNew)
	if err != nil {
		slog.Error("couldn't get newService key ", "err", err)
		return
	}
	newRes := resource.NewResource(newService.Name, newService.APIVersion, "", "service", newKey, portNameNew, newService)
	oldRes := resource.NewResource(oldService.Name, oldService.APIVersion, "", "service", oldKey, portNameOld, oldService)
	si.handler.OnUpdateService(newRes, oldRes)
}

func (si *ServiceInformer) OnDelete(obj interface{}) {
	service, ok := obj.(*v1.Service)
	if !ok {
		slog.Error("type of object is not service ", "obj", obj, "method", "OnDelete")
		return
	}

	portName, ok := isXDSService(service)
	if !ok {
		return
	}
	key, err := getServiceKey(service, portName)
	if err != nil {
		slog.Error("couldn't get service key ", "err", err)
		return
	}

	si.handler.DeleteService(key)
}
