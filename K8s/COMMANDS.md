#### CREATE CLUSTER

```bash
kind create cluster --config cluster.yaml
```

#### APPLY ALL CONFIGS

```bash
kubectl apply -k .
```

#### DESTROY CLUSTER

```bash
kind delete cluster 
```

#### DELETE ALL CONFIGS

```bash
kubectl delete -k .
```