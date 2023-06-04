package xds

import (
	"context"
	"net"
	"strconv"
	"time"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	xds "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	k8scache "k8s.io/client-go/tools/cache"
)

type ControlPlane struct {
	log               *logrus.Logger
	version           int32
	snapshotCache     cache.SnapshotCache
	server            xds.Server
	callBacks         *callbacks
	endpoints         []types.Resource
	endpointInformers []k8scache.SharedIndexInformer
	serviceInformers  []k8scache.SharedIndexInformer
	conf              *Config
	storage           cache.Storage
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

	informerEndpoints.AddEventHandler(k8scache.ResourceEventHandlerFuncs{
		UpdateFunc: cp.HandleEndpointsUpdate,
	})

	informerServices.AddEventHandler(k8scache.ResourceEventHandlerFuncs{
		UpdateFunc: cp.HandleServicesUpdate,
	})

	// go func() {
	// 	informerEndpoints.Run(stop)
	// }()

	// go func() {
	// 	informerServices.Run(stop)
	// }()

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
