package xds

import (
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
	"k8s.io/client-go/informers"
	k8scache "k8s.io/client-go/tools/cache"
	"net"
	"time"
)

type ControlPlane struct {
	log               *logrus.Logger
	version           int
	snapshotCache     cache.SnapshotCache
	server            xds.Server
	callBacks         *callbacks
	endpoints         []types.Resource
	endpointInformers []k8scache.SharedIndexInformer
}

func (cp *ControlPlane) Run() error {
	grpcServer := grpc.NewServer()

	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, cp.server)
	endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, cp.server)
	clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, cp.server)
	routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, cp.server)
	listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, cp.server)

	clusters, _ := CreateBootstrapClients()

	for _, cluster := range clusters {
		stop := make(chan struct{})
		defer close(stop)
		factory := informers.NewSharedInformerFactoryWithOptions(cluster, time.Second*10, informers.WithNamespace("demo"))
		informer := factory.Core().V1().Endpoints().Informer()
		cp.endpointInformers = append(cp.endpointInformers, informer)

		informer.AddEventHandler(k8scache.ResourceEventHandlerFuncs{
			UpdateFunc: cp.HandleEndpointsUpdate,
		})

		go func() {
			informer.Run(stop)
		}()
	}

	lis, _ := net.Listen("tcp", ":8080")
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}
