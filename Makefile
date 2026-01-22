# Catalyst Development Makefile
# Standardized workflow for local development

.PHONY: dev clean ui mock install help core test cluster-up cluster-down cluster-destroy prune deep-clean

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
	@echo "  make cluster-up      - Spin up Kind K8s Cluster (DB + Broker)"
	@echo "  make cluster-down    - Uninstall Helm Chart (Keep Cluster)"
	@echo "  make cluster-destroy - Delete Kind Cluster (Reclaim Resources)"
	@echo "  make prune           - Aggressive Docker Cleanup"
	@echo "  make deep-clean      - Nuke everything (Local + Cluster + Docker)"

# Nuke everything
deep-clean: clean cluster-destroy prune

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

sdk-check:
	@echo "Verifying SDK Usage..."
	@go list -f '{{.ImportPath}}: {{.Imports}}' ./... | findstr "github.com/point-unknown/catalyst/pkg" > NUL || echo "WARNING: Some packages may not be using the SDK"

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
	-kind create cluster --name catalyst-local
	helm install catalyst ./deploy/charts/catalyst -n catalyst-local --create-namespace

# Tear down Cluster
cluster-down:
	helm uninstall catalyst -n catalyst-local

# Destroy Kind Cluster (Reclaim Resources)
cluster-destroy:
	kind delete cluster --name catalyst-local

build-images:
	docker build -t catalyst/core:latest -f core/Dockerfile .
	docker build -t catalyst/ui:latest -f ui/Dockerfile ui

load-images: build-images
	kind load docker-image catalyst/core:latest --name catalyst-local
	kind load docker-image catalyst/ui:latest --name catalyst-local

# Aggressive Docker Cleanup (Save Disk Space)
prune:
	docker system prune -af --volumes

