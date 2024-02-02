#### CREATE CLUSTER

```bash
kind create cluster --config cluster.yaml
```

#### APPLY ALL CONFIGS

```bash
kubectl apply -k .
```

#### APPLY STATEFULSET

```bash
kubectl apply -k StatefulSet
```

#### APPLY DEPLOYMENT

```bash
kubectl apply -k Deployment
```
#### APPLY INGRESS

```bash
kubectl apply -k Ingress
```
#### DESTROY CLUSTER

```bash
kind delete cluster 
```

#### DELETE ALL CONFIGS

```bash
kubectl delete -k .
```

#### DELETE STATEFULSET

```bash
kubectl delete -k StatefulSet
```

#### DELETE DEPLOYMENT

```bash
kubectl delete -k Deployment
```