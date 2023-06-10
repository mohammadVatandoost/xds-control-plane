package xds

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	xds "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	k8scache "k8s.io/client-go/tools/cache"
)

type ControlPlane struct {
	log           *logrus.Logger
	version       int32
	snapshotCache cache.SnapshotCache
	server        xds.Server
	fetches       int32
	requests      int32
	// callBacks         *callbacks
	// endpoints         []types.Resource
	endpointInformers []k8scache.SharedIndexInformer
	serviceInformers  []k8scache.SharedIndexInformer
	conf              *Config
	storage           cache.Storage
	nodes             map[string]*Node
	mu                sync.RWMutex
	resources         map[string]map[string]struct{} // A resource is watched by which nodes
	muResource        sync.RWMutex
}

func (cp *ControlPlane) CreateNode(id string) *Node {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	node, ok := cp.nodes[id]
	if !ok {
		node = &Node{
			watchers: make(map[string]struct{}),
		}
	}
	cp.nodes[id] = node
	return node
}

func (cp *ControlPlane) GetNode(id string) (*Node, error) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	node, ok := cp.nodes[id]
	if !ok {
		return nil, fmt.Errorf("node with id: %s is not exist", id)
	}
	return node, nil
}

func (cp *ControlPlane) DeleteNode(id string) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	_, ok := cp.nodes[id]
	if !ok {
		return fmt.Errorf("node with id: %s is not exist", id)
	}
	delete(cp.nodes, id)
	return nil
}

func (cp *ControlPlane) AddResourceWatchToNode(id string, resource string) {
	cp.muResource.Lock()
	defer cp.muResource.Unlock()
	nodes, ok := cp.resources[resource]
	if !ok {
		nodes = make(map[string]struct{})
		cp.resources[resource] = nodes
	}
	nodes[id] = struct{}{}
}

func (cp *ControlPlane) GetNodesWatchTheResource(resource string) []string {
	cp.muResource.RLock()
	defer cp.muResource.RUnlock()
	nodesArray := make([]string, 0)
	nodes, ok := cp.resources[resource]
	if !ok {
		return nodesArray
	}
	for n := range nodes {
		nodesArray = append(nodesArray, n)
	}
	return nodesArray
}

func (cp *ControlPlane) Run() error {
	grpcServer := grpc.NewServer()

	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, cp.server)
	endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, cp.server)
	clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, cp.server)
	routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, cp.server)
	listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, cp.server)

	// clusters, _ := CreateBootstrapClients()
	clusterClient, err := CreateClusterClient()
	if err != nil {
		cp.log.WithError(err).Error("can not create cluster client")
		return err
	}
	cp.log.Info("cluster client created")
	namespaces, err := clusterClient.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{})
	if err != nil {
		cp.log.WithError(err).Error("can not get namespaces list")
		return err
	}
	cp.log.Infof("cluster number of namespaces: %v", len(namespaces.Items))
	for _, namespace := range namespaces.Items {
		cp.log.Infof("namespace: %v", namespace.Name)
	}
	cp.log.Info("==========")
	stop := make(chan struct{})
	defer close(stop)

	factory := informers.NewSharedInformerFactoryWithOptions(clusterClient, time.Second*10, informers.WithNamespace(""))

	informerEndpoints := factory.Core().V1().Endpoints().Informer()
	cp.endpointInformers = append(cp.endpointInformers, informerEndpoints)

	informerServices := factory.Core().V1().Services().Informer()
	cp.serviceInformers = append(cp.endpointInformers, informerServices)

	// informerEndpoints.AddEventHandler(k8scache.ResourceEventHandlerFuncs{
	// 	UpdateFunc: cp.HandleEndpointsUpdate,
	// })

	informerServices.AddEventHandler(k8scache.ResourceEventHandlerFuncs{
		UpdateFunc: cp.HandleServicesUpdate,
	})

	// go func() {
	// 	informerEndpoints.Run(stop)
	// }()

	go func() {
		informerServices.Run(stop)
	}()

	// go cp.RunXDSserver(stop)

	// for _, cluster := range clusters {
	// 	stop := make(chan struct{})
	// 	defer close(stop)
	// 	factory := informers.NewSharedInformerFactoryWithOptions(cluster, time.Second*10, informers.WithNamespace("demo"))
	// 	informer := factory.Core().V1().Endpoints().Informer()
	// 	cp.endpointInformers = append(cp.endpointInformers, informer)

	// 	informer.AddEventHandler(k8scache.ResourceEventHandlerFuncs{
	// 		UpdateFunc: cp.HandleEndpointsUpdate,
	// 	})

	// 	go func() {
	// 		informer.Run(stop)
	// 	}()
	// }

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(cp.conf.ListenPort))
	if err != nil {
		return err
	}

	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}
