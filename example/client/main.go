package main

import (
	"net"
	"strings"

	"github.com/mohammadVatandoost/interfaces/golang/echo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"google.golang.org/grpc/admin"
	_ "google.golang.org/grpc/resolver" // use for "dns:///be.cluster.local:50051"
	_ "google.golang.org/grpc/xds"      // use for xds-experimental:///be-srv
)

var (
	conn *grpc.ClientConn
)

type Config struct {
	Server1Address string
}

func main() {
	viper.SetDefault("Server1Address", "xds:///xds-grpc-server-example-headless:8888")
	// Read Config from ENV
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	var config Config

	err := viper.Unmarshal(&config)
	if err != nil {
		logrus.Fatalf("failed to read configs: %v", err)
	}


	//address = fmt.Sprintf("xds-experimental:///be-srv")

	// (optional) start background grpc admin services to monitor client
	// "google.golang.org/grpc/admin"
	go func() {
		lis, err := net.Listen("tcp", ":19000")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		defer lis.Close()
		opts := []grpc.ServerOption{grpc.MaxConcurrentStreams(10)}
		grpcServer := grpc.NewServer(opts...)
		cleanup, err := admin.Register(grpcServer)
		if err != nil {
			logrus.Fatalf("failed to register admin services: %v", err)
		}
		defer cleanup()

		logrus.Printf("GRPC Admin port listen on :%s", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			logrus.Fatalf("failed to serve: %v", err)
		}
	}()

	logrus.Printf("Connectting to server: %v ", config.Server1Address)

	conn, err := grpc.Dial(config.Server1Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := echo.NewEchoServerClient(conn)
	ctx := context.Background()
	i := 0
	for {
		r, err := c.SayHello(ctx, &echo.EchoRequest{Name: "unary RPC msg "})
		if err != nil {
			logrus.Fatalf("could not greet: %v", err)
		}
		logrus.Printf("RPC Response: %v %v", i, r)
		time.Sleep(5 * time.Second)
		i++
	}

}
