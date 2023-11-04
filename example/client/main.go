package main

import (
	"net"
	"net/http"
	"strings"

	"github.com/mohammadVatandoost/interfaces/golang/echo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"google.golang.org/grpc/admin"
	"google.golang.org/grpc/credentials/insecure"
	// "google.golang.org/grpc/credentials/xds"
	_ "google.golang.org/grpc/resolver" // use for "dns:///be.cluster.local:50051"
	_ "google.golang.org/grpc/xds"      // use for xds-experimental:///be-srv
)

var (
	conn *grpc.ClientConn
)

type Config struct {
	Server1Address string
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "client_rpc_call_counter",
		Help: "The total number of rpc calls",
	})
)

func main() {

	// viper.SetDefault("Server1Address", "xds:///xds-grpc-server-example-headless:8888")
	viper.SetDefault("Server1Address", "xds:///xds-grpc-server-example-headless.control-plane-example")
	// viper.SetDefault("Server1Address", "xds-grpc-server-example-headless:8888")
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

	logger := logrus.WithField("name", "client")
	var config Config

	err := viper.Unmarshal(&config)
	if err != nil {
		logger.Fatalf("failed to read configs: %v", err)
	}

	//address = fmt.Sprintf("xds-experimental:///be-srv")

	// (optional) start background grpc admin services to monitor client
	// "google.golang.org/grpc/admin"
	go func() {
		logger.Printf("Starting GRPC admin ")
		lis, err := net.Listen("tcp", ":19000")
		if err != nil {
			logger.Fatalf("failed to listen: %v", err)
		}
		defer lis.Close()
		opts := []grpc.ServerOption{grpc.MaxConcurrentStreams(10)}
		grpcServer := grpc.NewServer(opts...)
		cleanup, err := admin.Register(grpcServer)
		if err != nil {
			logger.Fatalf("failed to register admin services: %v", err)
		}
		defer cleanup()

		logger.Printf("GRPC Admin port listen on :%s", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			logrus.Fatalf("failed to serve: %v", err)
		}
	}()

	logger.Printf("Connectting to server: %v ", config.Server1Address)
	// creds, err := xds.NewClientCredentials(xds.ClientOptions{
	// 	FallbackCreds: insecure.NewCredentials(),
	// })

	conn, err := grpc.Dial(config.Server1Address,
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`), // This sets the initial balancing policy.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":9000", nil)
		if err != nil {
			logrus.WithError(err).Error("can not listen to expose metrics")
		}
	}()

	c := echo.NewEchoServerClient(conn)
	ctx := context.Background()
	i := 0
	for {
		r, err := c.SayHello(ctx, &echo.EchoRequest{Name: "unary RPC msg "})
		if err != nil {
			logrus.Errorf("could not greet: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		logger.Printf("RPC Response: %v %v", i, r)
		opsProcessed.Inc()
		time.Sleep(5 * time.Second)
		i++
	}

}
