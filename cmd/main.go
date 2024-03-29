package main

import (
	"context"
	"log/slog"
	"os"

	controlplane "github.com/mohammadVatandoost/xds-conrol-plane/internal/app/control-plane"
	"github.com/mohammadVatandoost/xds-conrol-plane/internal/core/rest"
	"github.com/mohammadVatandoost/xds-conrol-plane/internal/informer"
	"github.com/mohammadVatandoost/xds-conrol-plane/internal/k8s"
	"github.com/mohammadVatandoost/xds-conrol-plane/internal/xds"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/config"
	controlplaneConfig "github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/app/controlplane"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/utils"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/version"
)

const serviceName = "xds_control_plane"

func main() {
	slog.Info("Initializing", "service", serviceName, "information", version.Build.FormatDetailedProductInfo())
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	serverContext, serverCancel := utils.WithSignalCancellation(context.Background())
	defer serverCancel()

	conf := controlplaneConfig.DefaultControlPlaneConfig()
	err := config.Load("", conf)
	if err != nil {
		slog.Error("couldn't load configs", "error", err)
		exitCode = -1
		return
	}

	k8sClient, err := k8s.CreateClusterClient()
	if err != nil {
		slog.Error("couldn't create k8s client", "error", err)
		exitCode = -1
		return
	}

	cache := xds.NewSnapshotCache(conf.XDSConfig.ADSEnabled)

	app := controlplane.NewApp(conf, cache)

	restAPIServer := rest.NewServer(conf.RestAPIConfig, app)
	go func() {
		err := restAPIServer.Run()
		if err != nil {
			slog.Error("couldn't run rest API server", "error", err)
			exitCode = -1
		}

	}()

	runTimeInformer := informer.NewRunTime(k8sClient)
	serviceInformer := informer.NewServiceInformer(runTimeInformer.GetInformerFactory(), app)
	runTimeInformer.AddInformer(serviceInformer)
	runTimeInformer.RunInformers(serverContext.Done())

	slog.Info("XDS control plane config", "XDS.ADSEnabled", conf.XDSConfig.ADSEnabled, "ListenPort", conf.XDSConfig.Port)
	xdsServer := xds.NewControlPlane(conf.XDSConfig, app, cache)
	err = xdsServer.Run()
	if err != nil {
		slog.Error("couldn't run xdsServer", "error", err)
		exitCode = -1
	}
}
