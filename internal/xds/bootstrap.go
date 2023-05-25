package xds

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
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
	serviceName := svcc.ServiceName         //"be-srv"
	namespace := svcc.Namespace             //"default"
	portName := svcc.PortName               //"grpc"
	protocol := svcc.Protocol               //"tcp"
	grpcServiceName := svcc.GRPCServiceName //"echo.EchoServer"

	logrus.Printf("Looking up svc")
	cname, rec, err := net.LookupSRV(portName, protocol, fmt.Sprintf("%s.%s.svc.cluster.local", serviceName, namespace))
	if err != nil {
		logrus.Errorf("Could not find server %s", serviceName, err.Error())
		return upstreamPorts
	} else {
		logrus.Infof("SRV CNAME: %v, rec: %v\n", cname, rec)
	}

	var wg sync.WaitGroup

	for i := range rec {
		wg.Add(1)
		go func(host string, port string) {
			defer wg.Done()
			address := fmt.Sprintf("%s:%s", host, port)

			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, 30*time.Millisecond)
			defer cancel()
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				logrus.Errorf("Could not connect to endpoint %s  %v", address, err.Error())
				return
			}
			resp, err := healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{Service: grpcServiceName})
			if err != nil {
				logrus.WithField("address", address).Errorf("HealthCheck failed err: %v, conn: %v", conn, err.Error())
				// return // ToDo: for testign disable this
			}
			if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
				logrus.Errorf("Service not healthy %v %v", conn, fmt.Sprintf("service not in serving state: %v", resp.GetStatus().String()))
				// return ToDo: for testign disable this
			}
			logrus.Infof("RPC HealthChekStatus: for %v %v", address, resp.GetStatus())
			upstreamPorts = append(upstreamPorts, address)
		}(rec[i].Target, strconv.Itoa(int(rec[i].Port)))
	}
	wg.Wait()
	return upstreamPorts
}

func RunXDSserver(stopCh <-chan struct{}, snapshotCache cachev3.SnapshotCache) {
	var version int32
	for {
		select {
		case <-stopCh:
			return
		default:
			snapshot, err := makeSnapshot(version)
			if err != nil {
				logrus.Printf(">>>>>>>>>>  Error setting snapshot %v", err)
				return
			}
			IDs := snapshotCache.GetStatusKeys()
			logrus.Infof("snapshotCache IDs: %v\n", IDs)
			for _, id := range IDs {
				err = snapshotCache.SetSnapshot(context.Background(), id, snapshot)
				if err != nil {
					logrus.Errorf("%v", err)
				}
			}
		}
		time.Sleep(time.Duration(10) * time.Second)
	}
}

func makeSnapshot(version int32) (*cachev3.Snapshot, error) {
	services := getServices()
	eds := []types.Resource{}
	cls := []types.Resource{}
	rds := []types.Resource{}
	lsnr := []types.Resource{}
	for _, svc := range services {
		routeConfigName := svc.ServiceName + "-route"
		clusterName := svc.ServiceName + "-cluster"
		virtualHostName := svc.ServiceName + "-vs"
		region := svc.Region //"us-central1"
		zone := svc.Zone     // us-central1-a
		addresses := getAddresses(svc)
		logrus.Infof("service: %v, addresses: %v \n", svc.ServiceName, addresses)
		lbe := makeLBEndpoint(addresses)
		eds = append(eds, &endpoint.ClusterLoadAssignment{
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
		})
		cls = append(cls, &cluster.Cluster{
			Name:                 clusterName,
			LbPolicy:             cluster.Cluster_ROUND_ROBIN,
			ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_EDS},
			EdsClusterConfig: &cluster.Cluster_EdsClusterConfig{
				EdsConfig: &core.ConfigSource{
					ConfigSourceSpecifier: &core.ConfigSource_Ads{},
				},
			},
		})

		// RDS
		logrus.Infof(">>>>>>>>>>>>>>>>>>> creating RDS " + virtualHostName)
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

		rds = append(rds, &route.RouteConfiguration{
			Name:         routeConfigName,
			VirtualHosts: []*route.VirtualHost{vh},
		})

		// LISTENER
		logrus.Infof(">>>>>>>>>>>>>>>>>>> creating LISTENER " + svc.ServiceName)
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

		manager := &hcm.HttpConnectionManager{
			CodecType:      hcm.HttpConnectionManager_AUTO,
			RouteSpecifier: hcRds,
		}

		pbst, err := ptypes.MarshalAny(manager)
		if err != nil {
			panic(err)
		}

		lsnr = append(lsnr, &listener.Listener{
			Name: svc.ServiceName,
			ApiListener: &listener.ApiListener{
				ApiListener: pbst,
			},
		})
	}

	atomic.AddInt32(&version, 1)
	logrus.Infof(" creating snapshot Version " + fmt.Sprint(version))

	logrus.Infof("   snapshot with Listener %v", lsnr)
	logrus.Infof("   snapshot with EDS %v", eds)
	logrus.Infof("   snapshot with CLS %v", cls)
	logrus.Infof("   snapshot with RDS %v", rds)

	return cachev3.NewSnapshot(fmt.Sprint(version), map[resource.Type][]types.Resource{
		resource.EndpointType: eds,
		resource.ClusterType:  cls,
		resource.ListenerType: lsnr,
		resource.RouteType:    rds,
	})

	// snap := cachev3.NewSnapshot(fmt.Sprint(version), eds, cls, rds, lsnr, rt, sec)
	// err = config.SetSnapshot(nodeId, snap)

}

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
