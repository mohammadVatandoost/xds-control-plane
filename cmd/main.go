package main

import (
	"log/slog"
	"os"

	controlplane "github.com/mohammadVatandoost/xds-conrol-plane/internal/app/control-plane"
	"github.com/mohammadVatandoost/xds-conrol-plane/internal/informer"
	"github.com/mohammadVatandoost/xds-conrol-plane/internal/k8s"
	"github.com/mohammadVatandoost/xds-conrol-plane/internal/xds"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/config"
	controlplaneConfig "github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/app/controlplane"
)

const serviceName = "xds_control_plane"

func main() {
	exitCode := 0
	defer func ()  {
		os.Exit(exitCode)
	}()

	// ToDo: add signal cancellation
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

	app := controlplane.NewApp(conf)

	runTimeInformer := informer.NewRunTime(k8sClient)
	serviceInformer := informer.NewServiceInformer(runTimeInformer.GetInformerFactory(), app)
	runTimeInformer.AddInformer(serviceInformer)
	// runTimeInformer.RunInformers()

	slog.Info("XDS control plane config", "XDS.ADSEnabled", conf.XDSConfig.ADSEnabled, "ListenPort", conf.XDSConfig.Port)
	xdsServer := xds.NewControlPlane(conf.XDSConfig)
	err = xdsServer.Run()
	if err != nil {
		slog.Error("couldn't run xdsServer", "error", err)
		exitCode = -1
	}
}

