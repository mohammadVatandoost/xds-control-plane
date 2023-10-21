package informer

import (
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	k8scache "k8s.io/client-go/tools/cache"
)

type RunTime struct {
	client kubernetes.Interface
}


func (rt *RunTime) AddInformer(informer k8scache.ResourceEventHandler) {
	
	// informerServices.AddEventHandler(k8scache.ResourceEventHandlerFuncs{
	// 	UpdateFunc: cp.HandleServicesUpdate,
	// })
}


func (rt *RunTime) GetInformerFactory() informers.SharedInformerFactory {
	return informers.NewSharedInformerFactoryWithOptions(rt.client, time.Second*10, informers.WithNamespace(""))
}