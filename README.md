# Catalyst
> **The Autonomous Open Source Maintenance Platform.**

> 🚧 **Catalyst is currently in PRE-ALPHA (Architectural Incubation).**

## 🌍 Global Architecture State
**Current Strategy:** "Phase 1.5: Platform Dashboard"
**Goal:** Expand the "Vertical Slice" into a Multi-Source Control Plane, visualizing Hardware, Agent Activity, and Repository Events.

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
    Agents -->|agent/log| Broker
    
    Broker -->|WebSockets| Widgets
```

---

## 🤖 The Delivery Swarm
We utilize a virtual "Swarm" of specialized agent personas to execute this project.
**See `/docs/swarms.md` for current assignments.**

### 1. `UIPrime` (The Frontend Architect)
*   **Mission**: Deliver a "Cyber-Minimalist", premium-feel control plane.
*   **Tech Stack**: React 19, Vite, Tailwind CSS, Shadcn/UI, Framer Motion.
*   **Directives**:
    *   "Visualize everything; text is a fallback."
    *   "Sidebar navigation must be intuitive."
    *   **Responsibility**: `/ui` directory.

### 2. `SystemCore` (The Backend Engineer)
*   **Mission**: Build a fault-tolerant, high-concurrency event bus and orchestrator.
*   **Tech Stack**: Go (Golang), Eclipse Mosquitto (MQTT), PostgreSQL.
*   **Directives**:
    *   "The Event Bus is the source of truth."
    *   **Responsibility**: `/core` directory, `/bin/device-mock`.

### 3. `BridgeBuilder` (The Hardware Integrator)
*   **Mission**: Pierce the veil between the Container World and the Physical World.
*   **Tech Stack**: Go (Windows Native), `go-serial`, Win32 APIs.
*   **Directives**:
    *   "Hardware is messy; sanitize the inputs."
    *   **Responsibility**: `/bridge` directory.

---

## 🧠 Phase 2: Core Orchestration Architecture
The `catalyst-core` service is the central nervous system, built on a **Hexagonal Architecture** to ensure extensibility.

### 1. The Hexagonal Core
*   **Domain Layer** (`internal/domain`): [x] Contracts Defined (`CloudEvent`, `Agent`).
*   **Adapters** (`internal/adapter`):
    *   **EventBus**: [x] MQTT Client (Paho) Connected.
    *   **Store**: [ ] Persistence for active workflows.
*   **Service Layer** (`internal/service`):
    *   **MissionManager**: [x] Routing Logic Verified.
    *   **AgentRegistry**: [x] Plugin Loading (Concurrency Safe).

### 2. The Agent Execution Spectrum
Agents interact in 3 modes, rigorously typed in the Domain:
1.  **Reporting**: Telemetry/Logs (`agent.log`). Fire-and-forget. Verified via `ConsoleReporter`.
2.  **Communicating**: Inter-agent Signals (`agent.signal`). Verified via `TrendScout`.
3.  **Expressing**: Structured Artifacts (`data.report`). Final output.

### 3. Quality Assurance
*   **Unit Tests**: Comprehensive Go test suite for Domain, Service, and Agents.
    *   Run tests: `make test`
    *   Coverage: Event Marshaling, Registry Concurrency, Mission Routing, Anomaly Detection.

### 4. Extensibility
*   **New Inputs**: Core subscribes to wildcard topics (`sensor/#`).
*   **New Agents**: Implements the `Agent` interface (`Execute(ctx, event)`).

---

## 🛡️ Secure Self-Hosted Architecture (The "Ultrathink")

Catalyst is architected for **Zero Trust Local** execution.

*   **Isolation**: Services run in strictly isolated containers (Docker).
*   **Ingress Control**: No direct port exposure. Traffic flows through a Reverse Proxy.
*   **Least Privilege**:
    *   **Frontend**: Static Nginx build (No Node.js runtime in prod).
    *   **Backend**: Distroless Go binary.
    *   **Broker**: Mosquitto with explicit ACLs.

### 🛠️ Development Workflow
We adhere to a standardized `Makefile` workflow to ensure environment consistency.

**Prerequisites**: Docker Desktop, Go 1.22+, Node.js 20+.

```bash
# 1. Install Dependencies
make install

# 2. Start Full Stack (Default Ports)
# Launches UI (localhost:5173) and Device Mock
make dev

# 3. Clean Environment (Kill conflicting processes)
make clean

# 4. Start Individual Components
make ui    # Frontend only
make mock  # Backend Mock only
```

---

## 🛠️ Development Workflow (Local-First)
We prioritize **Localhost execution** over containerization during the "Incubation" phase.

### Phase 1.5: Platform Dashboard (Current)
1.  **Start Broker**: `docker compose up -d mosquitto` (Ports: `1883`, `9001`)
2.  **Start UI**: `cd ui && npm run dev` (Localhost: `5174`)
3.  **Start Mock System**: `go run ./bin/device-mock` (Simulates Hardware + Agents + Repo)

### Phase 2: K8s Deployment (Future)
*   Containerize `/core` and `/ui`.
*   Deploy to Kind cluster.
*   `BridgeBuilder` remains as a native Windows binary.

---

## 📂 Repository Structure
*   `/ui`: The React Frontend Application.
*   `/core`: The central Go Orchestrator.
*   `/bridge`: The Windows Native Hardware Bridge.
*   `/bin`: Helper scripts and mocks (e.g., `device-mock`).
*   `/deploy`: Docker Compose and Helm Charts.
*   `/docs`: Architecture Decision Records (ADRs).
