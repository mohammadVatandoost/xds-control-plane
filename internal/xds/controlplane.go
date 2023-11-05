package xds

import (
	"context"
	"log/slog"
	"net"
	"strconv"
	"sync"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	xds "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/mohammadVatandoost/xds-conrol-plane/internal/k8s"
	"github.com/mohammadVatandoost/xds-conrol-plane/internal/node"
	xdsConfig "github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/xds"
	"google.golang.org/grpc"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ControlPlane struct {
	version    int32
	cache      *XDSSnapshotCache
	server     xds.Server
	fetches    int32
	requests   int32
	conf       *xdsConfig.XDSConfig
	nodes      map[string]*node.Node
	mu         sync.RWMutex
	resources  map[string]map[string]struct{} // A resource is watched by which nodes
	muResource sync.RWMutex
	app        App
}

func (cp *ControlPlane) Run() error {
	grpcServer := grpc.NewServer()

	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, cp.server)
	endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, cp.server)
	clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, cp.server)
	routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, cp.server)
	listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, cp.server)

	// clusters, _ := CreateBootstrapClients()
	clusterClient, err := k8s.CreateClusterClient()
	if err != nil {
		slog.Error("can not create cluster client", "error", err)
		return err
	}
	slog.Info("cluster client created")
	namespaces, err := clusterClient.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{})
	if err != nil {
		slog.Error("can not get namespaces list", "error", err)
		return err
	}
	slog.Info("cluster", "NamespacesNum", len(namespaces.Items))
	for _, namespace := range namespaces.Items {
		slog.Info("", "namespace", namespace.Name)
	}
	slog.Info("==========")
	stop := make(chan struct{})
	defer close(stop)

	// factory := informers.NewSharedInformerFactoryWithOptions(clusterClient, time.Second*10, informers.WithNamespace(""))

	// informerEndpoints := factory.Core().V1().Endpoints().Informer()
	// cp.endpointInformers = append(cp.endpointInformers, informerEndpoints)

	// informerServices := factory.Core().V1().Services().Informer()
	// cp.serviceInformers = append(cp.endpointInformers, informerServices)

	// informerEndpoints.AddEventHandler(k8scache.ResourceEventHandlerFuncs{
	// 	UpdateFunc: cp.HandleEndpointsUpdate,
	// })

	// informerServices.AddEventHandler(k8scache.ResourceEventHandlerFuncs{
	// 	UpdateFunc: cp.HandleServicesUpdate,
	// })

	// go func() {
	// 	informerEndpoints.Run(stop)
	// }()

	// go func() {
	// 	informerServices.Run(stop)
	// }()

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(int(cp.conf.Port)))
	if err != nil {
		return err
	}

	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}
