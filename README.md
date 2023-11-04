# XDS Control Plane
A proxy-less service mesh for grpc services in kubernetes.

### Running in K8S

Use argoCD yaml files or Helm charts to deploy on K8s

### Running Locally by Kind
Setup dev tools
```shell
make dev/tools
```

Setup local k8s
```shell
make kind/start
```

Deploy xds-control-plane with server and client example servoce to k8s
```shell
make kind/deploy/control-plane
```

### ToDo:
- generate bootstrap file with tls (look at this https://github.com/mohammadVatandoost/traffic-director-grpc-bootstrap)
- reconcile GAMMA resources to config the traffic
- export metrics
