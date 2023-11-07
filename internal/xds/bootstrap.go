package xds

import (
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	v3routerpb "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"

	// "github.com/envoyproxy/go-control-plane/pkg/cache/types"
	// cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	// "github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/sirupsen/logrus"
	kube_core "k8s.io/api/core/v1"

	// "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	// healthpb "google.golang.org/grpc/health/grpc_health_v1"
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

// key is servicename.namespace
func getAddresses(key string, portName string) []string {
	var upstreamPorts []string
	// serviceName := svcc.ServiceName //"be-srv"
	// namespace := svcc.Namespace     //"default"
	// portName := svcc.PortName       //"grpc"
	protocol := kube_core.ProtocolTCP //"tcp"
	// grpcServiceName := svcc.GRPCServiceName //"echo.EchoServer"
	// name := fmt.Sprintf("%s.%s.svc.cluster.local", serviceName, namespace)
	// name := fmt.Sprintf("%s.%s.svc.cluster.local", serviceName, namespace)
	// name := fmt.Sprintf("%s.svc.cluster.local", key)
	name := key
	cname, rec, err := net.LookupSRV(portName, string(protocol), name)
	if err != nil {
		slog.Error("Could not find the address", "key", key, "portName", portName, "err", err.Error())
		return upstreamPorts
	} else {
		slog.Info("addresses found", "cname", cname, "rec", rec, "key", key, "name", name)
	}

	for i := range rec {
		address := fmt.Sprintf("%s:%s", rec[i].Target, strconv.Itoa(int(rec[i].Port)))
		upstreamPorts = append(upstreamPorts, address)
	}

	return upstreamPorts
}

// func (cp *ControlPlane) RunXDSserver(stopCh <-chan struct{}) {
// 	var version int32
// 	for {
// 		select {
// 		case <-stopCh:
// 			return
// 		default:
// 			snapshot, err := cp.makeSnapshot(version)
// 			if err != nil {
// 				slog.Error(">>>>>>>>>>  Error setting snapshot", "error", err)
// 				return
// 			}
// 			IDs := cp.snapshotCache.GetStatusKeys()
// 			slog.Info("snapshotCache", "IDs", IDs)
// 			for _, id := range IDs {
// 				err = cp.snapshotCache.SetSnapshot(context.Background(), id, snapshot)
// 				if err != nil {
// 					logrus.Errorf("%v", err)
// 				}
// 			}
// 		}
// 		time.Sleep(time.Duration(10) * time.Second)
// 	}
// }

// func (cp *ControlPlane) makeSnapshot(version int32) (*cachev3.Snapshot, error) {
// 	services := getServices()
// 	eds := []types.Resource{}
// 	cls := []types.Resource{}
// 	rds := []types.Resource{}
// 	lsnr := []types.Resource{}
// 	for _, svc := range services {
// 		edsService, clsService, rdsService, lsnrService, err := cp.makeXDSConfigFromService(svc)
// 		if err != nil {
// 			continue
// 		}
// 		eds = append(eds, edsService)
// 		cls = append(cls, clsService)
// 		rds = append(rds, rdsService)
// 		lsnr = append(lsnr, lsnrService)
// 	}

// 	atomic.AddInt32(&version, 1)
// 	slog.Info(" creating snapshot Version ", "version", version)

// 	slog.Info("   snapshot", "listner", lsnr)
// 	slog.Info("   snapshot with EDS %v", eds)
// 	slog.Info("   snapshot with CLS %v", cls)
// 	slog.Info("   snapshot with RDS %v", rds)

// 	return cachev3.NewSnapshot(fmt.Sprint(version), map[resource.Type][]types.Resource{
// 		resource.EndpointType: eds,
// 		resource.ClusterType:  cls,
// 		resource.ListenerType: lsnr,
// 		resource.RouteType:    rds,
// 	})
// }

func (cp *ControlPlane) makeXDSConfigFromService(svc ServiceConfig) (*endpoint.ClusterLoadAssignment, *cluster.Cluster, *route.RouteConfiguration, *listener.Listener, error) {
	routeConfigName := svc.ServiceName + "-route"
	clusterName := svc.ServiceName + "-cluster"
	virtualHostName := svc.ServiceName + "-vs"
	region := svc.Region //"us-central1"
	zone := svc.Zone     // us-central1-a
	addresses := getAddresses(fmt.Sprintf("%s.%s", svc.ServiceName, svc.PortName), svc.PortName)
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

	filterPbst, err := anypb.New(&v3routerpb.Router{})
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

	pbst, err := anypb.New(manager)
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
