# XDS Control Plane


## Running Locally

```shell


export GRPC_XDS_BOOTSTRAP= "./example/client/xds_bootstrap_local.json"
export GRPC_GO_LOG_VERBOSITY_LEVEL=99
export GRPC_GO_LOG_SEVERITY_LEVEL="info"
go run ./example/client/main.go
```

## Issues
[] for ADS, the request names must match the snapshot names, if they do not, then the watch is never responded, and it is expected that envoy makes another request. So we can only add service names to the snapshot that client exactly watch. this is wierld.
(WARN[0010] ADS mode: not responding to request: "kube-prometheus-prometheus:9090" not listed, ResourceNames: [xds-grpc-server-example-headless:8888] )