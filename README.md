# XDS Control Plane


## Running Locally

```shell


export GRPC_XDS_BOOTSTRAP= "./example/client/xds_bootstrap_local.json"
export GRPC_GO_LOG_VERBOSITY_LEVEL=99
export GRPC_GO_LOG_SEVERITY_LEVEL="info"
go run ./example/client/main.go
```