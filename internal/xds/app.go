package xds


type App interface {
	NewStreamRequest(id string, resourceNames []string) 
	StreamClosed(id string)
}