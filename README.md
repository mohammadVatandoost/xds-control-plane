# XDS Control Plane
A proxy-less service mesh for grpc services in kubernetes.
### Running in K8S

Use argoCD yaml files or Helm charts to deploy on K8s

### Running Locally by Kind
Setup local k8s
```shell
make kind/start
```
Deploy xds-control-plane with server and client example servoce to k8s
```shell
make kind/deploy/control-plane
```


### Issues
- [] for ADS, the request names must match the snapshot names, if they do not, then the watch is never responded, and it is expected that envoy makes another request. So we can only add service names to the snapshot that client exactly watch. this is wierld. It means if client watch xds-grpc-server-example-headless resource, you can only send listner for this resource (you couldn't resolve all the k8s services)
(WARN[0010] ADS mode: not responding to request: "kube-prometheus-prometheus:9090" not listed, ResourceNames: [xds-grpc-server-example-headless:8888] )
- [] Updating snapshot cache on XDS callbacks, cause a lot of requests send from client to XDS server ( Why????)
- [] With each informer update that is not neccessary, It update snapshot and send redunct update to XDS clients.

### ToDo:
- generate bootstrap file with tls (look at this https://github.com/mohammadVatandoost/traffic-director-grpc-bootstrap)
