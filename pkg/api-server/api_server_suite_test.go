package api_server_test

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"net"
// 	"net/http"
// 	"path/filepath"
// 	"testing"
// 	"time"

// 	. "github.com/onsi/gomega"

// 	mesh_proto "github.com/mohammadVatandoost/xds-conrol-plane/api/mesh/v1alpha1"
// 	api_server "github.com/mohammadVatandoost/xds-conrol-plane/pkg/api-server"
// 	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/api-server/customization"
// 	config_api_server "github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/api-server"
// 	kuma_cp "github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/app/kuma-cp"
// 	config_core "github.com/mohammadVatandoost/xds-conrol-plane/pkg/config/core"
// 	config_manager "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/config/manager"
// 	resources_access "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/access"
// 	core_mesh "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/apis/mesh"
// 	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/manager"
// 	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/model"
// 	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/model/rest"
// 	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/registry"
// 	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/resources/store"
// 	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/runtime"
// 	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/dns/vips"
// 	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/envoy/admin/access"
// 	core_metrics "github.com/mohammadVatandoost/xds-conrol-plane/pkg/metrics"
// 	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/authn/api-server/certs"
// 	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/resources/memory"
// 	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/test"
// 	test_runtime "github.com/mohammadVatandoost/xds-conrol-plane/pkg/test/runtime"
// 	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/tokens/builtin"
// 	xds_context "github.com/mohammadVatandoost/xds-conrol-plane/pkg/xds/context"
// 	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/xds/server"
// )

// func TestWs(t *testing.T) {
// 	test.RunSpecs(t, "API Server")
// }

// type resourceApiClient struct {
// 	address string
// 	path    string
// }

// func (r *resourceApiClient) fullAddress() string {
// 	return "http://" + r.address + r.path
// }

// func (r *resourceApiClient) get(name string) *http.Response {
// 	response, err := http.Get(r.fullAddress() + "/" + name)
// 	Expect(err).NotTo(HaveOccurred())
// 	return response
// }

// func (r *resourceApiClient) list() *http.Response {
// 	response, err := http.Get(r.fullAddress())
// 	Expect(err).NotTo(HaveOccurred())
// 	return response
// }

// func (r *resourceApiClient) delete(name string) *http.Response {
// 	request, err := http.NewRequest(
// 		"DELETE",
// 		r.fullAddress()+"/"+name,
// 		nil,
// 	)
// 	Expect(err).ToNot(HaveOccurred())
// 	response, err := http.DefaultClient.Do(request)
// 	Expect(err).ToNot(HaveOccurred())
// 	return response
// }

// func (r *resourceApiClient) put(res rest.Resource) *http.Response {
// 	jsonBytes, err := json.Marshal(res)
// 	Expect(err).ToNot(HaveOccurred())
// 	return r.putJson(res.GetMeta().Name, jsonBytes)
// }

// func (r *resourceApiClient) putJson(name string, json []byte) *http.Response {
// 	request, err := http.NewRequest(
// 		"PUT",
// 		r.fullAddress()+"/"+name,
// 		bytes.NewBuffer(json),
// 	)
// 	Expect(err).ToNot(HaveOccurred())
// 	request.Header.Add("content-type", "application/json")
// 	response, err := http.DefaultClient.Do(request)
// 	Expect(err).ToNot(HaveOccurred())
// 	return response
// }

// func putSampleResourceIntoStore(resourceStore store.ResourceStore, name string, mesh string) {
// 	resource := core_mesh.TrafficRouteResource{
// 		Spec: &mesh_proto.TrafficRoute{
// 			Conf: &mesh_proto.TrafficRoute_Conf{
// 				Destination: map[string]string{
// 					"path": "/sample-path",
// 				},
// 			},
// 		},
// 	}
// 	err := resourceStore.Create(context.Background(), &resource, store.CreateByKey(name, mesh))
// 	Expect(err).NotTo(HaveOccurred())
// }

// type testApiServerConfigurer struct {
// 	store   store.ResourceStore
// 	config  *config_api_server.ApiServerConfig
// 	metrics func() core_metrics.Metrics
// 	zone    string
// 	global  bool
// }

// func NewTestApiServerConfigurer() *testApiServerConfigurer {
// 	t := &testApiServerConfigurer{
// 		metrics: func() core_metrics.Metrics {
// 			m, _ := core_metrics.NewMetrics("Standalone")
// 			return m
// 		},
// 		config: config_api_server.DefaultApiServerConfig(),
// 		store:  memory.NewStore(),
// 	}
// 	t.config.GUI.Enabled = false
// 	return t
// }

// func (t *testApiServerConfigurer) WithZone(z string) *testApiServerConfigurer {
// 	t.global = false
// 	t.zone = z
// 	return t
// }

// func (t *testApiServerConfigurer) WithGlobal() *testApiServerConfigurer {
// 	t.zone = ""
// 	t.global = true
// 	return t
// }

// func (t *testApiServerConfigurer) WithStandalone() *testApiServerConfigurer {
// 	t.zone = ""
// 	t.global = false
// 	return t
// }

// func (t *testApiServerConfigurer) WithStore(resourceStore store.ResourceStore) *testApiServerConfigurer {
// 	t.store = resourceStore
// 	return t
// }

// // WithMetrics a function that creates metrics (needs to be a function as these can't be reused in case of failed startups)
// func (t *testApiServerConfigurer) WithMetrics(metricsFn func() core_metrics.Metrics) *testApiServerConfigurer {
// 	t.metrics = metricsFn
// 	return t
// }

// func (t *testApiServerConfigurer) WithConfigMutator(fn func(*config_api_server.ApiServerConfig)) *testApiServerConfigurer {
// 	fn(t.config)
// 	return t
// }

// func StartApiServer(t *testApiServerConfigurer) (*api_server.ApiServer, kuma_cp.Config, func()) {
// 	var apiServer *api_server.ApiServer
// 	var cfg kuma_cp.Config
// 	var stop func()

// 	Eventually(func() error {
// 		var err error
// 		apiServer, cfg, stop, err = tryStartApiServer(t)
// 		return err
// 	}).
// 		WithTimeout(time.Second * 30).
// 		WithPolling(time.Millisecond * 500).
// 		WithOffset(1).
// 		Should(Succeed())

// 	return apiServer, cfg, stop
// }

// func tryStartApiServer(t *testApiServerConfigurer) (*api_server.ApiServer, kuma_cp.Config, func(), error) {
// 	ctx, stop := context.WithCancel(context.Background())
// 	// we have to manually search for port and put it into config. There is no way to retrieve port of running
// 	// http.Server and we need it later for the client
// 	port, err := test.GetFreePort()
// 	if err != nil {
// 		return nil, kuma_cp.Config{}, stop, err
// 	}
// 	t.config.HTTP.Port = uint32(port)

// 	port, err = test.GetFreePort()
// 	if err != nil {
// 		return nil, kuma_cp.Config{}, stop, err
// 	}
// 	t.config.HTTPS.Port = uint32(port)
// 	if t.config.HTTPS.TlsKeyFile == "" {
// 		t.config.HTTPS.TlsKeyFile = filepath.Join("..", "..", "test", "certs", "server-key.pem")
// 		t.config.HTTPS.TlsCertFile = filepath.Join("..", "..", "test", "certs", "server-cert.pem")
// 		t.config.Auth.ClientCertsDir = filepath.Join("..", "..", "test", "certs", "client")
// 	}

// 	cfg := kuma_cp.DefaultConfig()
// 	cfg.ApiServer = t.config
// 	if t.zone != "" {
// 		cfg.Mode = config_core.Zone
// 		cfg.Multizone.Zone.Name = t.zone
// 	} else if t.global {
// 		cfg.Mode = config_core.Global
// 	}

// 	resManager := manager.NewResourceManager(t.store)
// 	apiServer, err := api_server.NewApiServer( //nolint:contextcheck
// 		resManager,
// 		xds_context.NewMeshContextBuilder(
// 			resManager,
// 			server.MeshResourceTypes(server.HashMeshExcludedResources),
// 			net.LookupIP,
// 			cfg.Multizone.Zone.Name,
// 			vips.NewPersistence(resManager, config_manager.NewConfigManager(t.store)),
// 			cfg.DNSServer.Domain,
// 			80,
// 		),
// 		customization.NewAPIList(),
// 		registry.Global().ObjectDescriptors(model.HasWsEnabled()),
// 		&cfg,
// 		t.metrics(),
// 		func() string { return "instance-id" },
// 		func() string { return "cluster-id" },
// 		certs.ClientCertAuthenticator,
// 		runtime.Access{
// 			ResourceAccess:       resources_access.NewAdminResourceAccess(cfg.Access.Static.AdminResources),
// 			DataplaneTokenAccess: nil,
// 			EnvoyAdminAccess: access.NewStaticEnvoyAdminAccess(
// 				cfg.Access.Static.ViewConfigDump,
// 				cfg.Access.Static.ViewStats,
// 				cfg.Access.Static.ViewClusters,
// 			),
// 		},
// 		&test_runtime.DummyEnvoyAdminClient{},
// 		builtin.TokenIssuers{
// 			DataplaneToken:   builtin.NewDataplaneTokenIssuer(resManager),
// 			ZoneIngressToken: builtin.NewZoneIngressTokenIssuer(resManager),
// 			ZoneToken:        builtin.NewZoneTokenIssuer(resManager),
// 		},
// 	)
// 	if err != nil {
// 		return nil, cfg, stop, err
// 	}
// 	errChan := make(chan error)
// 	go func() {
// 		err := apiServer.Start(ctx.Done()) //nolint:contextcheck
// 		errChan <- err
// 	}()

// 	tick := time.NewTicker(time.Millisecond * 500)
// 	leftTicks := 10
// 	for {
// 		if leftTicks == 0 {
// 			stop()
// 			return nil, cfg, stop, errors.New("no more ticks left")
// 		}
// 		select {
// 		case err = <-errChan:
// 			return nil, cfg, stop, err
// 		case <-tick.C:
// 			leftTicks--
// 			r, err := http.Get("http://" + apiServer.Address() + "/config")
// 			if err != nil {
// 				return nil, cfg, stop, err
// 			}
// 			r.Body.Close()
// 			if r.StatusCode == http.StatusOK {
// 				return apiServer, cfg, stop, nil
// 			}
// 		}
// 	}
// }
