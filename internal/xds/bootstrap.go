package xds

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	v3routerpb "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/util"
	"github.com/sirupsen/logrus"

	// "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

	// healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ServiceConfig struct {
	ServiceName     string
	Namespace       string
	PortName        string
	Protocol        string
	GRPCServiceName string
	Zone            string
	Region          string
}

func getServices() []ServiceConfig {
	return []ServiceConfig{
		{
			ServiceName:     "xds-grpc-server-example-headless",
			Namespace:       "test",
			PortName:        "grpc",
			GRPCServiceName: "echo.EchoServer",
			Protocol:        "tcp",
			Zone:            "us-central1-a",
			Region:          "us-central1",
		},
	}
}

func CreateClusterClient() (kubernetes.Interface, error) {
	homeDie, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	kubeConfigPath := homeDie + "/.kube/config"
	var config *rest.Config
	if utils.FileExists(kubeConfigPath) {
		logrus.Info("kube config file exist")
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, err
		}
	} else {
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	// dynamic.NewForConfig()
	return clientset, nil
}

func getAddresses(svcc ServiceConfig) []string {
	var upstreamPorts []string
	serviceName := svcc.ServiceName //"be-srv"
	namespace := svcc.Namespace     //"default"
	portName := svcc.PortName       //"grpc"
	protocol := svcc.Protocol       //"tcp"
	// grpcServiceName := svcc.GRPCServiceName //"echo.EchoServer"

	cname, rec, err := net.LookupSRV(portName, protocol, fmt.Sprintf("%s.%s.svc.cluster.local", serviceName, namespace))
	if err != nil {
		logrus.Errorf("Could not find serviceName: %s, portName: %v, protocol: %v,  err: %v", serviceName, portName, protocol, err.Error())
		return upstreamPorts
	} else {
		logrus.Infof("SRV CNAME: %v, rec: %v\n", cname, rec)
	}

	// var wg sync.WaitGroup

	for i := range rec {
		// wg.Add(1)
		// go func(host string, port string) {
		// 	defer wg.Done()
		// 	address := fmt.Sprintf("%s:%s", host, port)

		// 	ctx := context.Background()
		// 	ctx, cancel := context.WithTimeout(ctx, 30*time.Millisecond)
		// 	defer cancel()
		// 	// ToDo: handle health check later
		// 	// conn, err := grpc.Dial(address, grpc.WithInsecure())
		// 	// if err != nil {
		// 	// 	logrus.Errorf("Could not connect to endpoint %s  %v", address, err.Error())
		// 	// 	return
		// 	// }
		// 	// resp, err := healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{Service: grpcServiceName})
		// 	// if err != nil {
		// 	// 	logrus.WithField("address", address).Errorf("HealthCheck failed err: %v, conn: %v", conn, err.Error())
		// 	// 	// return // ToDo: for testign disable this
		// 	// }
		// 	// if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		// 	// 	logrus.Errorf("Service not healthy %v %v", conn, fmt.Sprintf("service not in serving state: %v", resp.GetStatus().String()))
		// 	// 	// return ToDo: for testign disable this
		// 	// }
		// 	// logrus.Infof("RPC HealthChekStatus: for %v %v", address, resp.GetStatus())
		// 	// upstreamPorts = append(upstreamPorts, address)
		// }(rec[i].Target, strconv.Itoa(int(rec[i].Port)))
		address := fmt.Sprintf("%s:%s", rec[i].Target, strconv.Itoa(int(rec[i].Port)))
		upstreamPorts = append(upstreamPorts, address)
	}
	// wg.Wait()

	return upstreamPorts
}

func (cp *ControlPlane) RunXDSserver(stopCh <-chan struct{}) {
	var version int32
	for {
		select {
		case <-stopCh:
			return
		default:
			snapshot, err := cp.makeSnapshot(version)
			if err != nil {
				slog.Error(">>>>>>>>>>  Error setting snapshot", "error", err)
				return
			}
			IDs := cp.snapshotCache.GetStatusKeys()
			slog.Info("snapshotCache", "IDs", IDs)
			for _, id := range IDs {
				err = cp.snapshotCache.SetSnapshot(context.Background(), id, snapshot)
				if err != nil {
					logrus.Errorf("%v", err)
				}
			}
		}
		time.Sleep(time.Duration(10) * time.Second)
	}
}

func (cp *ControlPlane) makeSnapshot(version int32) (*cachev3.Snapshot, error) {
	services := getServices()
	eds := []types.Resource{}
	cls := []types.Resource{}
	rds := []types.Resource{}
	lsnr := []types.Resource{}
	for _, svc := range services {
		edsService, clsService, rdsService, lsnrService, err := cp.makeXDSConfigFromService(svc)
		if err != nil {
			continue
		}
		eds = append(eds, edsService)
		cls = append(cls, clsService)
		rds = append(rds, rdsService)
		lsnr = append(lsnr, lsnrService)
	}

	atomic.AddInt32(&version, 1)
	slog.Info(" creating snapshot Version ", "version", version)

	slog.Info("   snapshot", "listner", lsnr)
	slog.Info("   snapshot with EDS %v", eds)
	slog.Info("   snapshot with CLS %v", cls)
	slog.Info("   snapshot with RDS %v", rds)

	return cachev3.NewSnapshot(fmt.Sprint(version), map[resource.Type][]types.Resource{
		resource.EndpointType: eds,
		resource.ClusterType:  cls,
		resource.ListenerType: lsnr,
		resource.RouteType:    rds,
	})
}

func (cp *ControlPlane) makeXDSConfigFromService(svc ServiceConfig) (*endpoint.ClusterLoadAssignment, *cluster.Cluster, *route.RouteConfiguration, *listener.Listener, error) {
	routeConfigName := svc.ServiceName + "-route"
	clusterName := svc.ServiceName + "-cluster"
	virtualHostName := svc.ServiceName + "-vs"
	region := svc.Region //"us-central1"
	zone := svc.Zone     // us-central1-a
	addresses := getAddresses(svc)
	if len(addresses) == 0 {
		return nil, nil, nil, nil, fmt.Errorf("there is no availabe address for service: %v", svc.ServiceName)
	}
	// cp.log.Infof("service: %v, addresses: %v \n", svc.ServiceName, addresses)
	lbe := makeLBEndpoint(addresses)
	eds := &endpoint.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpoint.LocalityLbEndpoints{{
			Locality: &core.Locality{
				Region: region,
				Zone:   zone,
			},
			Priority:            0,
			LoadBalancingWeight: &wrapperspb.UInt32Value{Value: uint32(1000)},
			LbEndpoints:         lbe,
		}},
	}
	cls := &cluster.Cluster{
		Name:                 clusterName,
		LbPolicy:             cluster.Cluster_ROUND_ROBIN,
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_EDS},
		EdsClusterConfig: &cluster.Cluster_EdsClusterConfig{
			EdsConfig: &core.ConfigSource{
				ConfigSourceSpecifier: &core.ConfigSource_Ads{},
			},
		},
	}

	// RDS
	// cp.log.Infof(">>>>>>>>>>>>>>>>>>> creating RDS " + virtualHostName)
	vh := &route.VirtualHost{
		Name:    virtualHostName,
		Domains: []string{svc.ServiceName}, //******************* >> must match what is specified at xds:/// //

		Routes: []*route.Route{{
			Match: &route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Prefix{
					Prefix: "",
				},
			},
			Action: &route.Route_Route{
				Route: &route.RouteAction{
					ClusterSpecifier: &route.RouteAction_Cluster{
						Cluster: clusterName,
					},
				},
			},
		}}}

	rds := &route.RouteConfiguration{
		Name:         routeConfigName,
		VirtualHosts: []*route.VirtualHost{vh},
	}

	// LISTENER
	// cp.log.Infof(">>>>>>>>>>>>>>>>>>> creating LISTENER " + svc.ServiceName)
	hcRds := &hcm.HttpConnectionManager_Rds{
		Rds: &hcm.Rds{
			RouteConfigName: routeConfigName,
			ConfigSource: &core.ConfigSource{
				ConfigSourceSpecifier: &core.ConfigSource_Ads{
					Ads: &core.AggregatedConfigSource{},
				},
			},
		},
	}

	filterPbst, err := ptypes.MarshalAny(&v3routerpb.Router{})
	if err != nil {
		panic(err)
	}
	// RouterHTTPFilter := hcm.HTTPFilter("router", &v3routerpb.Router{})
	RouterHTTPFilter := &hcm.HttpFilter{
		Name: "router",
		ConfigType: &hcm.HttpFilter_TypedConfig{
			TypedConfig: filterPbst,
		},
	}
	filters := []*hcm.HttpFilter{
		RouterHTTPFilter,
	}

	manager := &hcm.HttpConnectionManager{
		CodecType:      hcm.HttpConnectionManager_AUTO,
		RouteSpecifier: hcRds,
		HttpFilters:    filters,
	}

	pbst, err := ptypes.MarshalAny(manager)
	if err != nil {
		panic(err)
	}

	lsnr := &listener.Listener{
		Name: svc.ServiceName,
		ApiListener: &listener.ApiListener{
			ApiListener: pbst,
		},
	}

	return eds, cls, rds, lsnr, nil
}

// // HTTPFilter constructs an xds HttpFilter with the provided name and config.
// func HTTPFilter(name string, config proto.Message) *hcm.HttpFilter {
// 	return &hcm.HttpFilter{
// 		Name: name,
// 		ConfigType: &hcm.HttpFilter_TypedConfig{
// 			TypedConfig: MarshalAny(config),
// 		},
// 	}
// }

// // MarshalAny is a convenience function to marshal protobuf messages into any
// // protos. It will panic if the marshaling fails.
// func MarshalAny(m proto.Message) *anypb.Any {
// 	a, err := ptypes.MarshalAny(m)
// 	if err != nil {
// 		panic(fmt.Sprintf("ptypes.MarshalAny(%+v) failed: %v", m, err))
// 	}
// 	return a
// }

func makeLBEndpoint(addresses []string) []*endpoint.LbEndpoint {
	lbe := make([]*endpoint.LbEndpoint, 0)
	for _, v := range addresses {
		backendHostName := strings.Split(v, ":")[0]
		backendPort := strings.Split(v, ":")[1]
		uPort, err := strconv.ParseUint(backendPort, 10, 32)
		if err != nil {
			logrus.Errorf("Could not parse port %v", err)
			break
		}
		// ENDPOINT
		logrus.Infof(">>>>>>>>>>>>>>>>>>> creating ENDPOINT for remoteHost:port %s:%s", backendHostName, backendPort)
		hst := &core.Address{Address: &core.Address_SocketAddress{
			SocketAddress: &core.SocketAddress{
				Address:  backendHostName,
				Protocol: core.SocketAddress_TCP,
				PortSpecifier: &core.SocketAddress_PortValue{
					PortValue: uint32(uPort),
				},
			},
		}}

		ee := &endpoint.LbEndpoint{
			HostIdentifier: &endpoint.LbEndpoint_Endpoint{
				Endpoint: &endpoint.Endpoint{
					Address: hst,
				}},
			HealthStatus: core.HealthStatus_HEALTHY,
		}
		lbe = append(lbe, ee)
	}
	return lbe
}
