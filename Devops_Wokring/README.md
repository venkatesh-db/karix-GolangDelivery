
Docker ::

docker-compose up

curl http://localhost:8080/health

 curl http://localhost:8082/order
{"id":1,"item":"Laptop","price":1200}

 curl http://localhost:8081/user
{"id":1,"name":"John Doe"}


venkatesh@venkateshs-MacBook-Pro ~ % curl http://localhost:8082/order
{"id":1,"item":"Laptop","price":1200}
venkatesh@venkateshs-MacBook-Pro ~ % curl http://localhost:8081/user
{"id":1,"name":"John Doe"}
venkatesh@venkateshs-MacBook-Pro ~ % curl http://localhost:8080/health
{"status":"ok","service":"api-gateway"}

Kuberntes :

 
Easy Working process  :

% kubectl port-forward svc/api-gateway -n demo-go-micro 8080:80

kubectl port-forward svc/user-service -n demo-go-micro 8081:8081

kubectl port-forward svc/order-service -n demo-go-micro 8082:8082

curl http://localhost:8081/user
curl http://localhost:8082/order



Troubleshooting process :

kubectl describe svc api-gateway -n demo-go-micro

kubectl get pods -l app=api-gateway -n demo-go-micro


kubectl logs -l app=order-service -n demo-go-micro

kubectl logs -l app=user-service -n demo-go-micro

kubectl describe svc user-service -n demo-go-micro

kubectl describe svc order-service -n demo-go-micro


kubectl port-forward svc/api-gateway -n demo-go-micro 8080:80

kubectl port-forward svc/order-service -n demo-go-micro 8082:80

kubectl port-forward svc/user-service -n demo-go-micro 8081:80

kubectl port-forward svc/api-gateway -n demo-go-micro 8080:80


curl http://localhost:8080/health
curl http://localhost:8082/order
curl http://localhost:8081/user

kubectl get events -n demo-go-micro





My laptop  process ::

kubectl delete pods --all --namespace=demo-go-micro


eval $(minikube docker-env) && docker build -t api-gateway:latest ./api-gateway && docker build -t user-service:latest ./user-service && docker build -t order-service:latest ./order-service


kubectl apply -f k8s/

kubectl get pods -n demo-go-micro

kubectl get svc -n demo-go-micro

 kubectl get pods -n demo-go-micro

kubectl logs -n demo-go-micro -l app=api-gateway
kubectl logs -n demo-go-micro -l app=order-service
kubectl logs -n demo-go-micro -l app=user-service

kubectl delete pod -n demo-go-micro api-gateway-5fd848b5fc-2gjhr


kubectl port-forward svc/api-gateway -n demo-go-micro 8080:80


kubectl exec -it <pod-name> -n demo-go-micro -- curl http://api-gateway.demo-go-micro.svc.cluster.local/health


minikube service api-gateway -n demo-go-micro

venkatesh@venkateshs-MacBook-Pro ~ % kubectl describe svc api-gateway -n demo-go-micro
Name:                     api-gateway
Namespace:                demo-go-micro
Labels:                   <none>
Annotations:              <none>
Selector:                 app=api-gateway
Type:                     NodePort
IP Family Policy:         SingleStack
IP Families:              IPv4
IP:                       10.107.241.180
IPs:                      10.107.241.180
Port:                     <unset>  80/TCP
TargetPort:               8080/TCP
NodePort:                 <unset>  30081/TCP
Endpoints:                10.244.0.243:8080,10.244.0.247:8080,10.244.0.249:8080
Session Affinity:         None
External Traffic Policy:  Cluster
Events:                   <none>


venkatesh@venkateshs-MacBook-Pro ~ % kubectl port-forward svc/api-gateway -n demo-go-micro 8080:80
Unable to listen on port 8080: Listeners failed to create with the following errors: [unable to create listener: Error listen tcp4 127.0.0.1:8080: bind: address already in use unable to create listener: Error listen tcp6 [::1]:8080: bind: address already in use]
error: unable to listen on any of the requested ports: [{8080 8080}]


# Go Microservices Project

This project demonstrates a complete example of a multi-service Go application with Docker, Docker Compose, and Kubernetes. It includes structured logging, Prometheus metrics, and horizontal scaling.

---

## **Project Structure**

```
.
├── api-gateway
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── user-service
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── order-service
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── docker-compose.yml
└── k8s
    ├── namespace.yaml
    ├── api-gateway-deployment.yaml
    ├── user-service-deployment.yaml
    ├── order-service-deployment.yaml
```

---

## **Services**

### 1. `api-gateway`
- **Port**: 8080
- **Endpoints**:
  - `GET /health`: Returns `{"status":"ok","service":"api-gateway"}`.
  - `GET /hello`: Calls `user-service` and `order-service` and returns combined JSON.
- **Metrics**: Exposed at `/metrics`.
- **Logging**: Structured logging with `log/slog`.

### 2. `user-service`
- **Port**: 8081
- **Endpoints**:
  - `GET /user`: Returns hardcoded user JSON.
- **Metrics**: Exposed at `/metrics`.
- **Logging**: Structured logging with `log/slog`.

### 3. `order-service`
- **Port**: 8082
- **Endpoints**:
  - `GET /order`: Returns hardcoded order JSON.
- **Metrics**: Exposed at `/metrics`.
- **Logging**: Structured logging with `log/slog`.

---

## **How to Run**

### 1. **Run with Docker Compose**

```bash
# Navigate to the project directory
cd /Users/venkatesh/Golang\ WOW\ Placments/devops/Devops

# Build the Docker images
docker-compose build

# Start the services
docker-compose up
```

- Access the `api-gateway` at [http://localhost:8080](http://localhost:8080).
- Scale the `api-gateway` service:
  ```bash
  docker-compose up --scale api-gateway=3
  ```

---

### 2. **Run with Kubernetes (Minikube)**

#### Step 1: Start Minikube
```bash
minikube start
```

#### Step 2: Build Docker Images in Minikube
```bash
eval $(minikube docker-env)

# Build images for all services
docker build -t api-gateway:latest ./api-gateway
docker build -t user-service:latest ./user-service
docker build -t order-service:latest ./order-service
```

#### Step 3: Apply Kubernetes Manifests
```bash
kubectl apply -f k8s/
```

#### Step 4: Access the Application
```bash
minikube service api-gateway -n demo-go-micro
```

---

## **Validation**

### Logs
Check logs for structured logging:
```bash
kubectl logs -f <pod-name> -n demo-go-micro
```

### Metrics
Access Prometheus metrics at `/metrics` for each service.

---

## **Features**
- **Structured Logging**: Using `log/slog`.
- **Prometheus Metrics**: HTTP request count and latency.
- **Horizontal Scaling**: Demonstrated in both Docker Compose and Kubernetes.

---

## **Kubernetes Details**

### Namespace
- `demo-go-micro`

### Deployments
- `api-gateway`: 3 replicas, NodePort service.
- `user-service`: 2 replicas, ClusterIP service.
- `order-service`: 2 replicas, ClusterIP service.

### Scaling and Load Balancing
- Kubernetes Service load balances between replicas.
- Verify scaling with:
  ```bash
  kubectl get pods -n demo-go-micro
  ```

---

## **License**
This project is licensed under the MIT License.