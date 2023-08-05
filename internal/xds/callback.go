package xds

import (
	"context"
	"sync/atomic"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
)

func (cp *ControlPlane) Report() {
	log.Info("cp.Report()  callbacks", "fetches", cp.fetches, "requests", cp.requests)
}
func (cp *ControlPlane) OnStreamOpen(ctx context.Context, id int64, typ string) error {
	log.Info("OnStreamOpen", "id", id, "typ", typ)
	return nil
}
func (cp *ControlPlane) OnStreamClosed(id int64, node *core.Node) {
	log.Info("OnStreamClosed", "id", id, "node.Id", node.Id, "node.Cluster", node.Cluster)
	cp.DeleteNode(node.Id)
}
func (cp *ControlPlane) OnStreamRequest(id int64, r *discovery.DiscoveryRequest) error {
	if r.TypeUrl != resource.ListenerType {
		return nil
	}
	log.Info("OnStreamRequest ", "id", id, "TypeUrl", r.TypeUrl, "ResourceNames", r.ResourceNames)
	node := cp.CreateNode(r.Node.Id)
	for _, rn := range r.ResourceNames {
		node.AddWatcher(rn)
		cp.AddResourceWatchToNode(r.Node.Id, rn)
	}
	return nil
}
func (cp *ControlPlane) OnStreamResponse(ctx context.Context, id int64, req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {
	log.Info("OnStreamResponse... %d   Request [%v],  Response[%v]", id, req.TypeUrl, resp.TypeUrl)
	cp.Report()
}
func (cp *ControlPlane) OnFetchRequest(ctx context.Context, req *discovery.DiscoveryRequest) error {
	log.Info("OnFetchRequest... Request [%v]", req.TypeUrl)
	atomic.AddInt32(&cp.fetches, 1)
	return nil
}
func (cp *ControlPlane) OnFetchResponse(req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {
	log.Info("OnFetchResponse... Resquest[%v],  Response[%v]", req.TypeUrl, resp.TypeUrl)
}

func (cp *ControlPlane) OnDeltaStreamClosed(id int64, node *core.Node) {
	log.Info("OnDeltaStreamClosed... %v", id)
}

func (cp *ControlPlane) OnDeltaStreamOpen(ctx context.Context, id int64, typ string) error {
	log.Info("OnDeltaStreamOpen... %v  of type %s", id, typ)
	return nil
}

func (cp *ControlPlane) OnStreamDeltaRequest(i int64, request *discovery.DeltaDiscoveryRequest) error {
	log.Info("OnStreamDeltaRequest... %v  of type %s", i, request)
	return nil
}

func (cp *ControlPlane) OnStreamDeltaResponse(i int64, request *discovery.DeltaDiscoveryRequest, response *discovery.DeltaDiscoveryResponse) {
	log.Info("OnStreamDeltaResponse... %v  of type %s", i, request)
}
