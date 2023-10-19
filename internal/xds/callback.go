package xds

import (
	"context"
	"log/slog"
	"sync/atomic"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
)

func (cp *ControlPlane) Report() {
	slog.Info("cp.Report()  callbacks", "fetches", cp.fetches, "requests", cp.requests)
}
func (cp *ControlPlane) OnStreamOpen(ctx context.Context, id int64, typ string) error {
	slog.Info("OnStreamOpen open for Type", "id", id, "type", typ)
	return nil
}
func (cp *ControlPlane) OnStreamClosed(id int64, node *core.Node) {
	slog.Info("OnStreamClosed closed", "id", id)
	cp.DeleteNode(node.Id)
}
func (cp *ControlPlane) OnStreamRequest(id int64, r *discovery.DiscoveryRequest) error {
	if r.TypeUrl != resource.ListenerType {
		return nil
	}
	slog.Info("OnStreamRequest", "id", id, "request", r.TypeUrl, "ResourceNames", r.ResourceNames)
	node := cp.CreateNode(r.Node.Id)
	for _, rn := range r.ResourceNames {
		node.AddWatcher(rn)
		cp.AddResourceWatchToNode(r.Node.Id, rn)
	}
	return nil
}
func (cp *ControlPlane) OnStreamResponse(ctx context.Context, id int64, req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {
	slog.Info("OnStreamResponse Request,  Response", "id", id, "request", req.TypeUrl, "response", resp.TypeUrl)
	cp.Report()
}
func (cp *ControlPlane) OnFetchRequest(ctx context.Context, req *discovery.DiscoveryRequest) error {
	slog.Info("OnFetchRequest...", "request", req.TypeUrl)
	atomic.AddInt32(&cp.fetches, 1)
	return nil
}
func (cp *ControlPlane) OnFetchResponse(req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {
	slog.Info("OnFetchResponse... Resquest[%v],  Response[%v]", "request", req.TypeUrl, "response", resp.TypeUrl)
}

func (cp *ControlPlane) OnDeltaStreamClosed(id int64, node *core.Node) {
	slog.Info("OnDeltaStreamClosed...", "id", id)
}

func (cp *ControlPlane) OnDeltaStreamOpen(ctx context.Context, id int64, typ string) error {
	slog.Info("OnDeltaStreamOpen...", "id", id, "type", typ)
	return nil
}

func (cp *ControlPlane) OnStreamDeltaRequest(i int64, request *discovery.DeltaDiscoveryRequest) error {
	slog.Info("OnStreamDeltaRequest... ", "id", i, "request", request)
	return nil
}

func (cp *ControlPlane) OnStreamDeltaResponse(i int64, request *discovery.DeltaDiscoveryRequest, response *discovery.DeltaDiscoveryResponse) {
	slog.Info("OnStreamDeltaResponse... ", "id", i, "request", request)
}
