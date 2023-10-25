package xds


type App interface {
	NewStreamRequest(id string, resourceNames []string, typeURL string) 
	StreamClosed(id string)
}