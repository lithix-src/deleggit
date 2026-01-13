# Deploy Script
kind create cluster --config deploy/k8s/kind-config.yaml --name catalyst-local
kubectl apply -f deploy/k8s/namespace.yaml
kubectl apply -f deploy/k8s/postgres.yaml
kubectl apply -f deploy/k8s/mosquitto.yaml
kubectl get pods -n catalyst-local -w
