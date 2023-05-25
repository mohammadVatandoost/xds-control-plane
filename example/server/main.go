package main

import (
	"net"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/mohammadVatandoost/interfaces/golang/echo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Config struct {
	GRPCPort   int
	ServerName string
}

var (
	hs   *health.Server
	conn *grpc.ClientConn
)

type server struct {
	ServerName string
}

func isGrpcRequest(r *http.Request) bool {
	return r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc")
}

func (s *server) SayHello(ctx context.Context, in *echo.EchoRequest) (*echo.EchoReply, error) {

	logrus.Println("Got rpc: --> ", in.Name)

	return &echo.EchoReply{Message: "Hello " + in.Name + "  from " + s.ServerName}, nil
}

func (s *server) SayHelloStream(in *echo.EchoRequest, stream echo.EchoServer_SayHelloStreamServer) error {

	logrus.Println("Got stream:  -->  ")
	stream.Send(&echo.EchoReply{Message: "Hello " + in.Name})
	stream.Send(&echo.EchoReply{Message: "Hello " + in.Name})

	return nil
}

type healthServer struct{}

func (s *healthServer) Check(ctx context.Context, in *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	logrus.Printf("Handling grpc Check request")
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (s *healthServer) Watch(in *healthpb.HealthCheckRequest, srv healthpb.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Watch is not implemented")
}

func main() {

	// Setting defaults for this application

	viper.SetDefault("GRPCPort", 8888)

	viper.SetDefault("ServerName", "server1")

	// Read Config from ENV
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "@timestamp",
			logrus.FieldKeyMsg:  "message",
		},
	})
	logrus.SetLevel(logrus.TraceLevel)

	var config Config

	err := viper.Unmarshal(&config)
	if err != nil {
		logrus.Fatalf("failed to read configs: %v", err)
	}

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(config.GRPCPort))
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	creds := insecure.NewCredentials()
	sopts := []grpc.ServerOption{grpc.MaxConcurrentStreams(10), grpc.Creds(creds)}
	// s := xds.NewGRPCServer(sopts...)
	s := grpc.NewServer(sopts...)

	echo.RegisterEchoServerServer(s, &server{ServerName: config.ServerName})

	healthpb.RegisterHealthServer(s, &healthServer{})
	logrus.Infof("Starting grpcServer on Port: %v with servername: %v", config.GRPCPort, config.ServerName)
	s.Serve(lis)

}
