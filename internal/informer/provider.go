package informer

import (
	"sync"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

type RunTime struct {
	client kubernetes.Interface
	informers []Informer
	mu sync.Mutex
}


func (rt *RunTime) AddInformer(informer Informer) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	rt.informers = append(rt.informers, informer)
}

func (rt *RunTime) RunInformers(stopCh <- chan struct{}) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	for _, informer := range rt.informers {
		go informer.Run(stopCh)
	}
}


func (rt *RunTime) GetInformerFactory() informers.SharedInformerFactory {
	return informers.NewSharedInformerFactoryWithOptions(rt.client, time.Second*10, informers.WithNamespace(""))
}


func NewRunTime(client kubernetes.Interface) *RunTime {
	return &RunTime{
		client: client,
	}
}