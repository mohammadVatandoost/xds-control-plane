package controlplane

import "encoding/json"

func (a *App) GetResources() ([]byte, error) {
	return json.Marshal(a.resources)
}

func (a *App) GetNodes() ([]byte, error) {
	return json.Marshal(a.nodes)
}
