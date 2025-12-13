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