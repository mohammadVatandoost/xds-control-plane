package informer


type Informer interface {
	Run(stopCh <-chan struct{})
}

