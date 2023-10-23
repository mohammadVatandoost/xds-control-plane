package controlplane

import "log/slog"


func (a *App) NewStreamRequest(id string, resourceNames []string) {
	node := a.CreateNode(id)
	for _, rn := range resourceNames {
		node.AddWatcher(rn)
		a.AddResourceWatchToNode(id, rn)
	}
}

func (a *App) StreamClosed(id string) {
	err := a.DeleteNode(id)
	if err != nil {
		slog.Error("app stream closed", "error", err, "id", id)
	}
}
