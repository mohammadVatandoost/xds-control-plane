
## Tis is for client secret

```shell
kubectl create secret generic client-secret --from-file=cert.pem --dry-run=true  --output=yaml > secret.yaml
kubectl apply -f secret.yaml -n control-plane
```

