kind create cluster --config K8s/cluster.yaml
kubectl apply -f K8s/zookeeper.yaml
kubectl apply -f K8s/kafka.yaml
kubectl apply -f K8s/minio.yaml
kubectl apply -f K8s/processing_server.yaml
kubectl apply -f K8s/rabbitmq.yaml
kubectl apply -f K8s/command_server.yaml
kubectl apply -f K8s/mongodb.yaml
kubectl apply -f K8s/auth_server.yaml
kubectl apply -f K8s/main_server.yaml
kubectl apply -f K8s/notification_bot.yaml
kubectl apply -f K8s/conversation_bot.yaml
kubectl apply -f K8s/prometheus.yaml
kubectl apply -f K8s/sla_manager.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
kubectl apply -f K8s/ingress.yaml