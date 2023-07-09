
## Tis is for kuma client secret

```shell
kubectl create secret generic client-secret-for-kuma --from-file=cert.pem --dry-run=true  --output=yaml > secret.yaml
kubectl apply -f secret.yaml -n kuma-system
```

