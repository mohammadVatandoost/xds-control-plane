package xds

import (
	"context"
	"sync/atomic"

	"github.com/sirupsen/logrus"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
)

// type callbacks struct {
// 	signal        chan struct{}
// 	fetches       int
// 	requests      int
// 	mu            sync.Mutex
// 	log           *logrus.Logger
// 	eventsHandler EventsHandler
// }

// type EventsHandler interface {
// 	AddR(nodeID string, resourceNames []string, resourceType string)
// }

func (cp *ControlPlane) Report() {
	cp.log.WithFields(logrus.Fields{"fetches": cp.fetches, "requests": cp.requests}).Info("cp.Report()  callbacks")
}
func (cp *ControlPlane) OnStreamOpen(ctx context.Context, id int64, typ string) error {
	cp.log.Infof("OnStreamOpen %d open for Type [%s]", id, typ)
	return nil
}
func (cp *ControlPlane) OnStreamClosed(id int64, node *core.Node) {
	cp.log.Infof("OnStreamClosed %d closed, node id: %v, node cluster: %v", id, node.Id, node.Cluster)
	cp.DeleteNode(node.Id)
}
func (cp *ControlPlane) OnStreamRequest(id int64, r *discovery.DiscoveryRequest) error {
	if r.TypeUrl != resource.ListenerType {
		return nil
	}
	cp.log.Infof("OnStreamRequest %d  Request[%v], ResourceNames: %v", id, r.TypeUrl, r.ResourceNames)
	node := cp.CreateNode(r.Node.Id)
	for _, r := range r.ResourceNames {
		node.AddWatcher(r)
	}
	// cp.eventsHandler.UpdateCache(r.Node.Id, r.ResourceNames, r.TypeUrl)
	return nil
}
func (cp *ControlPlane) OnStreamResponse(ctx context.Context, id int64, req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {
	cp.log.Infof("OnStreamResponse... %d   Request [%v],  Response[%v]", id, req.TypeUrl, resp.TypeUrl)
	cp.Report()
}
func (cp *ControlPlane) OnFetchRequest(ctx context.Context, req *discovery.DiscoveryRequest) error {
	cp.log.Infof("OnFetchRequest... Request [%v]", req.TypeUrl)
	atomic.AddInt32(&cp.fetches, 1)
	return nil
}
func (cp *ControlPlane) OnFetchResponse(req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {
	cp.log.Infof("OnFetchResponse... Resquest[%v],  Response[%v]", req.TypeUrl, resp.TypeUrl)
}

func (cp *ControlPlane) OnDeltaStreamClosed(id int64, node *core.Node) {
	cp.log.Infof("OnDeltaStreamClosed... %v", id)
}

func (cp *ControlPlane) OnDeltaStreamOpen(ctx context.Context, id int64, typ string) error {
	cp.log.Infof("OnDeltaStreamOpen... %v  of type %s", id, typ)
	return nil
}

func (cp *ControlPlane) OnStreamDeltaRequest(i int64, request *discovery.DeltaDiscoveryRequest) error {
	cp.log.Infof("OnStreamDeltaRequest... %v  of type %s", i, request)
	return nil
}

func (cp *ControlPlane) OnStreamDeltaResponse(i int64, request *discovery.DeltaDiscoveryRequest, response *discovery.DeltaDiscoveryResponse) {
	cp.log.Infof("OnStreamDeltaResponse... %v  of type %s", i, request)
}
