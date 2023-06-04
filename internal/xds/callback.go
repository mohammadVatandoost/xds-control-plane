package xds

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
)

type callbacks struct {
	signal        chan struct{}
	fetches       int
	requests      int
	mu            sync.Mutex
	log           *logrus.Logger
	eventsHandler EventsHandler
}

type EventsHandler interface {
	UpdateCache(nodeID string, resourceNames []string, resourceType string)
}

func (cb *callbacks) Report() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.log.WithFields(logrus.Fields{"fetches": cb.fetches, "requests": cb.requests}).Info("cb.Report()  callbacks")
}
func (cb *callbacks) OnStreamOpen(ctx context.Context, id int64, typ string) error {
	cb.log.Infof("OnStreamOpen %d open for Type [%s]", id, typ)
	return nil
}
func (cb *callbacks) OnStreamClosed(id int64, node *core.Node) {
	cb.log.Infof("OnStreamClosed %d closed, node id: %v, node cluster: %v", id, node.Id, node.Cluster)
}
func (cb *callbacks) OnStreamRequest(id int64, r *discovery.DiscoveryRequest) error {
	cb.log.Infof("OnStreamRequest %d  Request[%v], ResourceNames: %v", id, r.TypeUrl, r.ResourceNames)
	cb.eventsHandler.UpdateCache(r.Node.Id, r.ResourceNames, r.TypeUrl)
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.requests++
	if cb.signal != nil {
		close(cb.signal)
		cb.signal = nil
	}
	return nil
}
func (cb *callbacks) OnStreamResponse(ctx context.Context, id int64, req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {
	cb.log.Infof("OnStreamResponse... %d   Request [%v],  Response[%v]", id, req.TypeUrl, resp.TypeUrl)
	cb.Report()
}
func (cb *callbacks) OnFetchRequest(ctx context.Context, req *discovery.DiscoveryRequest) error {
	cb.log.Infof("OnFetchRequest... Request [%v]", req.TypeUrl)
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.fetches++
	if cb.signal != nil {
		close(cb.signal)
		cb.signal = nil
	}
	return nil
}
func (cb *callbacks) OnFetchResponse(req *discovery.DiscoveryRequest, resp *discovery.DiscoveryResponse) {
	cb.log.Infof("OnFetchResponse... Resquest[%v],  Response[%v]", req.TypeUrl, resp.TypeUrl)
}

func (cb *callbacks) OnDeltaStreamClosed(id int64, node *core.Node) {
	cb.log.Infof("OnDeltaStreamClosed... %v", id)
}

func (cb *callbacks) OnDeltaStreamOpen(ctx context.Context, id int64, typ string) error {
	cb.log.Infof("OnDeltaStreamOpen... %v  of type %s", id, typ)
	return nil
}

func (cb *callbacks) OnStreamDeltaRequest(i int64, request *discovery.DeltaDiscoveryRequest) error {
	cb.log.Infof("OnStreamDeltaRequest... %v  of type %s", i, request)
	return nil
}

func (cb *callbacks) OnStreamDeltaResponse(i int64, request *discovery.DeltaDiscoveryRequest, response *discovery.DeltaDiscoveryResponse) {
	cb.log.Infof("OnStreamDeltaResponse... %v  of type %s", i, request)
}
