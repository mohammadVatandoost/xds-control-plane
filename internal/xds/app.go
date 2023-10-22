package xds


type App interface {
	NewStreamRequest(id string) 
	StreamClosed(id string)
}