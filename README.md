# Catalyst
> **The Autonomous Open Source Maintenance Platform.**

> 🚀 **Catalyst MVP 1.0 (Swarm Activated)**

## 🌍 Global Architecture State
**Current Strategy:** "Phase 2: Core Orchestration & Swarm Intelligence"
**Goal:** A self-hosted Control Plane ("High-Visibility Industrial Slate") visualizing Hardware Telemetry, Agent Swarm Operations, and Repository Events.

### 🏗️ Design Reference (The "Architecture Constraint")
All functional components MUST adhere to this topology:

```mermaid
graph TD
    subgraph "Windows Host (Source Layer)"
        Device[Hardware Bridge (Go)]
        Repo[Repository Watcher (Go)]
        Agents[Agent Swarm (Go/Containers)]
    end

    subgraph "Catalyst Runtime (Event Bus)"
        Broker[Mosquitto MQTT]
    end

    subgraph "Platform Dashboard (UI Layer)"
        Sidebar[Nav: Dashboard/Workflows/Hardware]
        Widgets[Grid: Sensors/Logs/Events]
    end

    Device -->|sensor/cpu/temp| Broker
    Repo -->|repo/issue/new| Broker
    Agents -->|agent/+/log| Broker
    
    Broker -->|WebSockets| Widgets
```

---

## 🤖 The Delivery Swarm (MVP 1.0)
We utilize a virtual "Swarm" of specialized agent personas to execute this project.
**See `AGENTS.md` for detailed functional specs.**

### 1. `Interface` (The Frontend Architect)
*   **Role**: Frontend Engineering & UX.
*   **Mission**: Deliver a **High-Visibility "Industrial Slate"** control plane optimized for monitoring hardware.
*   **Tech Stack**: React 19, Vite, Tailwind CSS (Slate Palette), Framer Motion.
*   **Directives**: "Visualize everything. Reduce cognitive load."

### 2. `Orchestrator` (The Backend Engineer)
*   **Role**: Core Systems Engineering.
*   **Mission**: Build a fault-tolerant, high-concurrency event bus and orchestrator.
*   **Tech Stack**: Go (Golang), Eclipse Mosquitto (MQTT), PostgreSQL (Persistence).
*   **Directives**: "The Event Bus is the source of truth."

### 3. `Infrastructure` (DevOps & Site Reliability)
*   **Role**: Infrastructure-as-Code (IaC) & K8s.
*   **Mission**: Maintain the `kind` cluster, Docker resources, and self-hosted environments.
*   **Tech Stack**: Kubernetes, Helm, Docker, Cygwin/Make.

### 4. `Simulation` (Data Generation)
*   **Role**: Chaos Engineering.
*   **Mission**: Emulate hardware sensors and swarm activity for development (`device-mock`).
*   **Status**: Active (Emitting `agent/+/log` and `sensor/cpu/temp`).

---

## 🧠 Phase 2: Core Orchestration Architecture
The `catalyst-core` service is the central nervous system, built on a **Hexagonal Architecture**.

### 1. The Hexagonal Core
*   **Domain Layer** (`internal/domain`): [x] Contracts Defined (`CloudEvent`, `Agent`).
*   **Adapters** (`internal/adapter`):
    *   **EventBus**: [x] MQTT Client (Paho) Connected.
    *   **Store**: [x] Postgres Persistence (Event Logging).
*   **Service Layer** (`internal/service`):
    *   **MissionManager**: [x] Routing Logic Verified.
    *   **AgentRegistry**: [x] Plugin Loading (Concurrency Safe).

### 2. The Agent Execution Spectrum
Agents interact in 3 modes, rigorously typed in the Domain:
1.  **Reporting**: Telemetry/Logs (`agent.log`). Fire-and-forget.
2.  **Communicating**: Inter-agent Signals (`agent.signal`).
3.  **Expressing**: Structured Artifacts (`data.report`). Final output.

---

## 🛡️ Secure Self-Hosted Architecture

Catalyst is architected for **Zero Trust Local** execution.

*   **Isolation**: Services run in strictly isolated containers (Docker).
*   **Ingress Control**: No direct port exposure. Traffic flows through a Reverse Proxy.
*   **Least Privilege**: Strict ACLs on the MQTT Broker.

### 🛠️ Development Workflow
We adhere to a standardized `Makefile` workflow.

**Prerequisites**: Docker Desktop, Go 1.22+, Node.js 20+, Kind.

```bash
# 1. Install Dependencies & Tools
make install

# 2. Start Full Stack (Localhost)
# Launches UI (localhost:5173), Core Service, and Device Mock
make dev

# 3. Clean Environment (Kill processes, prune containers)
make clean

# 4. Infrastructure Management
make cluster-up   # Start Kind Cluster (DB/Broker)
make cluster-down # Destroy Cluster
```

---

## 📂 Repository Structure
*   `/ui`: The React Frontend Application (`Interface`).
*   `/core`: The central Go Orchestrator (`Orchestrator`).
*   `/bin`: Helper scripts and mocks (`Simulation`).
*   `/deploy`: Kubernetes Manifests and Docker configs (`Infrastructure`).
*   `/docs`: Architecture Decision Records (ADRs).
