package k8s

import (
	"context"
	"fmt"

	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/logger"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	corev1 "k8s.io/api/core/v1"
	
)

var log = logger.NewLoggerWithName("k8s-event-Listener")

type Emitter interface {
	OnCreate(kind string, key ResourceKey)
	OnUpdate(kind string, key ResourceKey)
	OnDelete(kind string, key ResourceKey)
}

type Listener struct {
	mgr manager.Manager
	out Emitter
}

func NewListener(mgr manager.Manager, out Emitter) *Listener {
	return &Listener{
		mgr: mgr,
		out: out,
	}
}

func (k *Listener) Start(stop <-chan struct{}) error {
	types := core_registry.Global().ObjectTypes()
	knownTypes := k.mgr.GetScheme().KnownTypes(kuma_v1alpha1.GroupVersion)
	for _, t := range types {
		if _, ok := knownTypes[string(t)]; !ok {
			continue
		}
		gvk := kuma_v1alpha1.GroupVersion.WithKind(string(t))
		lw, err := k.createListerWatcher(gvk)
		if err != nil {
			return err
		}
		coreObj, err := core_registry.Global().NewObject(t)
		if err != nil {
			return err
		}
		obj, err := k8s_registry.Global().NewObject(coreObj.GetSpec())
		if err != nil {
			return err
		}

		serviceType := corev1.Service{}
		informer := cache.NewSharedInformer(lw, serviceType, 0)
		if _, err := informer.AddEventHandler(k); err != nil {
			return err
		}

		go func(typ string) {
			log.V(1).Info("start watching resource", "type", typ)
			informer.Run(stop)
		}(t)
	}
	return nil
}

func resourceKey(obj KubernetesObject) ResourceKey {
	var name string
	if obj.Scope() == ScopeCluster {
		name = obj.GetName()
	} else {
		name = fmt.Sprintf("%s.%s", obj.GetName(), obj.GetNamespace())
	}
	return ResourceKey{
		Name: name,
		Mesh: obj.GetMesh(),
	}
}

func (k *Listener) OnAdd(obj interface{}, _ bool) {
	kobj := obj.(KubernetesObject)
	if err := k.addTypeInformationToObject(kobj); err != nil {
		log.Error(err, "unable to add TypeMeta to KubernetesObject")
		return
	}
	k.out.OnCreate(kobj.GetObjectKind().GroupVersionKind().Kind, resourceKey(kobj) )
}

func (k *Listener) OnUpdate(oldObj, newObj interface{}) {
	kobj := newObj.(KubernetesObject)
	if err := k.addTypeInformationToObject(kobj); err != nil {
		log.Error(err, "unable to add TypeMeta to KubernetesObject")
		return
	}
	k.out.OnUpdate(kobj.GetObjectKind().GroupVersionKind().Kind, resourceKey(kobj) )
}

func (k *Listener) OnDelete(obj interface{}) {
	kobj := obj.(KubernetesObject)
	if err := k.addTypeInformationToObject(kobj); err != nil {
		log.Error(err, "unable to add TypeMeta to KubernetesObject")
		return
	}
	k.out.OnCreate(kobj.GetObjectKind().GroupVersionKind().Kind, resourceKey(kobj) )
}

func (k *Listener) NeedLeaderElection() bool {
	return false
}

func (k *Listener) addTypeInformationToObject(obj runtime.Object) error {
	gvks, _, err := k.mgr.GetScheme().ObjectKinds(obj)
	if err != nil {
		return errors.Wrap(err, "missing apiVersion or kind and cannot assign it")
	}

	for _, gvk := range gvks {
		if len(gvk.Kind) == 0 {
			continue
		}
		if len(gvk.Version) == 0 || gvk.Version == runtime.APIVersionInternal {
			continue
		}
		obj.GetObjectKind().SetGroupVersionKind(gvk)
		break
	}

	return nil
}

func (k *Listener) createListerWatcher(gvk schema.GroupVersionKind) (cache.ListerWatcher, error) {
	mapping, err := k.mgr.GetRESTMapper().RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, err
	}
	httpClient, err := rest.HTTPClientFor(k.mgr.GetConfig())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HTTP client from Manager config")
	}
	client, err := apiutil.RESTClientForGVK(gvk, false, k.mgr.GetConfig(), serializer.NewCodecFactory(k.mgr.GetScheme()), httpClient)
	if err != nil {
		return nil, err
	}
	listGVK := gvk.GroupVersion().WithKind(gvk.Kind + "List")
	listObj, err := k.mgr.GetScheme().New(listGVK)
	if err != nil {
		return nil, err
	}
	paramCodec := runtime.NewParameterCodec(k.mgr.GetScheme())
	ctx := context.Background()
	return &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			res := listObj.DeepCopyObject()
			err := client.Get().
				Resource(mapping.Resource.Resource).
				VersionedParams(&opts, paramCodec).
				Do(ctx).
				Into(res)
			return res, err
		},
		// Setup the watch function
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			// Watch needs to be set to true separately
			opts.Watch = true
			return client.Get().
				Resource(mapping.Resource.Resource).
				VersionedParams(&opts, paramCodec).
				Watch(ctx)
		},
	}, nil
}
