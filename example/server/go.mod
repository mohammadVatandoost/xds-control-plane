module server

go 1.19

require (
	github.com/mohammadVatandoost/interfaces/golang v0.0.0-20230210155147-08d0f820f3fb
	golang.org/x/net v0.5.0
	google.golang.org/grpc v1.53.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace echo => ../echo
