kind delete cluster
kubectl delete -f K8s/zookeeper.yaml
kubectl delete -f K8s/kafka.yaml
kubectl delete -f K8s/minio.yaml
kubectl delete -f K8s/processing_server.yaml
kubectl delete -f K8s/rabbitmq.yaml
kubectl delete -f K8s/command_server.yaml
kubectl delete -f K8s/mongodb.yaml
kubectl delete -f K8s/auth_server.yaml
kubectl delete -f K8s/main_server.yaml
kubectl delete -f K8s/notification_bot.yaml
kubectl delete -f K8s/conversation_bot.yaml
kubectl delete -f K8s/prometheus.yaml
kubectl delete -f K8s/sla_manager.yaml
kubectl delete -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
kubectl delete -f K8s/ingress.yaml