# Catalyst Development Makefile
# Standardized workflow for local development

.PHONY: dev clean ui mock install help core test

help:
	@echo "Catalyst Dev Environment"
	@echo "  make dev      - Start the full stack (UI + Device Mock) on default ports"
	@echo "  make clean    - Kill all Catalyst-related processes and free ports"
	@echo "  make ui       - Start only the Frontend (localhost:5173)"
	@echo "  make mock     - Start only the Backend Mock"
	@echo "  make core     - Start the Core Service"
	@echo "  make test     - Run all Core Unit Tests"
	@echo "  make install  - Install all dependencies"
	@echo "  --- Infrastructure ---"
	@echo "  make cluster-up   - Spin up Kind K8s Cluster (DB + Broker)"
	@echo "  make cluster-down - Delete Kind Cluster"
	@echo "  make prune        - Aggressive Docker Cleanup"

# Install dependencies (One-time setup)
install:
	cd ui && npm install

# --- Standardization ---
lint: ## Enforce Engineering Standards (Go & TS)
	@echo "Checking Frontend Standards..."
	cd ui && npm run lint
	@echo "Checking Backend Standards..."
	go fmt ./...
	@echo "Standards verification complete."

# Clean up all ports and processes
clean:
	@echo "Cleaning up environment..."
	-npx -y kill-port 5173
	-npx -y kill-port 5174
	-npx -y kill-port 5175
	-npx -y kill-port 5176
	-npx -y kill-port 5177
	-npx -y kill-port 5178
	-taskkill /F /IM device-mock.exe 2>NUL || echo "Mock was not running"

# Start the full development stack
dev: clean build-mock build-watcher
	@echo "Starting Catalyst Stack..."
	@start "Catalyst Mock" cmd /c "bin\device-mock\device-mock.exe"
	@start "Repo Watcher" cmd /c "bin\repo-watcher\repo-watcher.exe"
	@cd ui && start "Catalyst UI" cmd /c "npm run dev"
	@echo "Stack launched! UI: http://localhost:5173"

# Build the Device Mock binary
build-mock:
	cd bin/device-mock && go build -o device-mock.exe main.go

# Build the Repo Watcher binary
build-watcher:
	cd bin/repo-watcher && go build -o repo-watcher.exe main.go

# Start only the UI
ui:
	cd ui && npm run dev

# Start only the Backend Mock
mock:
	cd bin/device-mock && go run main.go

# Start the Core Service
core:
	cd core && go run cmd/server/main.go

# Run Tests
test:
	cd core && go test ./... -v

# ------------------------------------
# Infrastructure (Kind & Docker)
# ------------------------------------

# Spin up Local Kubernetes Cluster (Postgres + Mosquitto + Observability)
cluster-up:
	kind create cluster --config deploy/k8s/kind-config.yaml --name catalyst-local
	kubectl apply -f deploy/k8s/namespace.yaml
	kubectl apply -f deploy/k8s/postgres.yaml
	kubectl apply -f deploy/k8s/mosquitto.yaml
	kubectl apply -f deploy/k8s/observability.yaml
	@echo "⏳ Waiting for Pods..."
	kubectl get pods -n catalyst-local -w


# Tear down Cluster
cluster-down:
	kind delete cluster --name catalyst-local

# Aggressive Docker Cleanup (Save Disk Space)
prune:
	@echo "🧹 Pruning unused Docker objects..."
	docker system prune -a -f --volumes

