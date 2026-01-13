# Deleggit Development Makefile
# Standardized workflow for local development

.PHONY: dev clean ui mock install help

help:
	@echo "Deleggit Dev Environment"
	@echo "  make dev      - Start the full stack (UI + Device Mock) on default ports"
	@echo "  make clean    - Kill all Deleggit-related processes and free ports"
	@echo "  make ui       - Start only the Frontend (localhost:5173)"
	@echo "  make mock     - Start only the Device Bridge Mock"
	@echo "  make install  - Install all dependencies"

# Install dependencies (One-time setup)
install:
	cd ui && npm install

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
dev: clean
	@echo "Starting Deleggit Stack..."
	@start "Deleggit Mock" cmd /c "bin\device-mock\device-mock.exe"
	@cd ui && start "Deleggit UI" cmd /c "npm run dev"
	@echo "Stack launched! UI: http://localhost:5173"

# Start only the UI
ui:
	cd ui && npm run dev

# Start only the Backend Mock
mock:
	bin\device-mock\device-mock.exe
