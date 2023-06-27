package customization_test

import (
	"net"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"

	api_server "github.com/mohammadVatandoost/xds-conrol-plane/pkg/api-server"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/api-server/customization"
	config_api_server "github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/api-server"
	kuma_cp "github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/app/kuma-cp"
	config_manager "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/config/manager"
	resources_access "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/access"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/manager"
	core_model "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/model"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/registry"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/store"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/runtime"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/dns/vips"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/envoy/admin/access"
	core_metrics "github.com/mohammadVatandoost/xds-conrol-plane/pkg/metrics"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/authn/api-server/certs"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/test"
	test_runtime "github.com/mohammadVatandoost/xds-conrol-plane/pkg/test/runtime"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/tokens/builtin"
	xds_context "github.com/mohammadVatandoost/xds-conrol-plane/pkg/xds/context"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/xds/server"
)

func TestWs(t *testing.T) {
	test.RunSpecs(t, "API Server Customization")
}

func createTestApiServer(store store.ResourceStore, config *config_api_server.ApiServerConfig, metrics core_metrics.Metrics, wsManager customization.APIManager) *api_server.ApiServer {
	// we have to manually search for port and put it into config. There is no way to retrieve port of running
	// http.Server and we need it later for the client
	port, err := test.GetFreePort()
	Expect(err).NotTo(HaveOccurred())
	config.HTTP.Port = uint32(port)

	port, err = test.GetFreePort()
	Expect(err).NotTo(HaveOccurred())
	config.HTTPS.Port = uint32(port)
	if config.HTTPS.TlsKeyFile == "" {
		config.HTTPS.TlsKeyFile = filepath.Join("..", "..", "..", "test", "certs", "server-key.pem")
		config.HTTPS.TlsCertFile = filepath.Join("..", "..", "..", "test", "certs", "server-cert.pem")
		config.Auth.ClientCertsDir = filepath.Join("..", "..", "..", "test", "certs", "client")
	}

	if wsManager == nil {
		wsManager = customization.NewAPIList()
	}
	cfg := kuma_cp.DefaultConfig()
	cfg.ApiServer = config
	resManager := manager.NewResourceManager(store)
	apiServer, err := api_server.NewApiServer(
		resManager,
		xds_context.NewMeshContextBuilder(
			resManager,
			server.MeshResourceTypes(server.HashMeshExcludedResources),
			net.LookupIP,
			cfg.Multizone.Zone.Name,
			vips.NewPersistence(resManager, config_manager.NewConfigManager(store)),
			cfg.DNSServer.Domain,
			cfg.DNSServer.ServiceVipPort,
		),
		wsManager,
		registry.Global().ObjectDescriptors(core_model.HasWsEnabled()),
		&cfg,
		metrics,
		func() string { return "instance-id" },
		func() string { return "cluster-id" },
		certs.ClientCertAuthenticator,
		runtime.Access{
			ResourceAccess:       resources_access.NewAdminResourceAccess(cfg.Access.Static.AdminResources),
			DataplaneTokenAccess: nil,
			EnvoyAdminAccess:     access.NoopEnvoyAdminAccess{},
		},
		&test_runtime.DummyEnvoyAdminClient{},
		builtin.TokenIssuers{
			DataplaneToken:   builtin.NewDataplaneTokenIssuer(resManager),
			ZoneIngressToken: builtin.NewZoneIngressTokenIssuer(resManager),
			ZoneToken:        builtin.NewZoneTokenIssuer(resManager),
		},
	)
	Expect(err).ToNot(HaveOccurred())
	return apiServer
}
