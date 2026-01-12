# Deleggit
> **The Autonomous Open Source Maintenance Platform.**

> 🚧 **Deleggit is currently in PRE-ALPHA (Architectural Incubation).**

## 🌍 Global Architecture State
**Current Strategy:** "The Vertical Slice MVP"
**Goal:** Establish a full data loop from Hardware to Web UI immediately, bypassing initial K8s complexity for rapid prototyping.

### 🏗️ Design Reference (The "Architecture Constraint")
All functional components MUST adhere to this topology:

```mermaid
graph TD
    subgraph "Windows Host (Hardware Layer)"
        Device[Hardware / USB / Serial]
        Bridge[Device Bridge (Go Binary)]
    end

    subgraph "Deleggit Runtime (Local/K8s)"
        Broker[MQTT Broker (Mosquitto)]
        Core[Deleggit Core (Go)]
        Agents[Agent Swarm (TrendScout, GapAnalyst...)]
    end

    subgraph "User Interface"
        UI[React + Shadcn UI (Cyber-Minimalist)]
    end

    Device <-->|Serial/HID| Bridge
    Bridge <-->|MQTT (TCP)| Broker
    Core <-->|MQTT (TCP)| Broker
    Agents <-->|MQTT (TCP)| Broker
    UI <-->|MQTT (WebSockets)| Broker
```

---

## 🤖 The Delivery Swarm
We utilize a virtual "Swarm" of specialized agent personas to execute this project.

### 1. `UIPrime` (The Frontend Architect)
*   **Mission**: Deliver a "Cyber-Minimalist", premium-feel control plane.
*   **Tech Stack**: React 19, Vite, Tailwind CSS, Shadcn/UI, Framer Motion.
*   **Directives**:
    *   "If it looks generic, it is wrong."
    *   "Visualize everything; text is a fallback."
    *   **Responsibility**: `/ui` directory.

### 2. `SystemCore` (The Backend Engineer)
*   **Mission**: Build a fault-tolerant, high-concurrency event bus and orchestrator.
*   **Tech Stack**: Go (Golang), Eclipse Mosquitto (MQTT), PostgreSQL.
*   **Directives**:
    *   "Concurrency is not parallelism."
    *   "The Event Bus is the source of truth."
    *   **Responsibility**: `/core` directory, `/bin/device-mock`.

### 3. `BridgeBuilder` (The Hardware Integrator)
*   **Mission**: Pierce the veil between the Container World and the Physical World.
*   **Tech Stack**: Go (Windows Native), `go-serial`, Win32 APIs.
*   **Directives**:
    *   "Hardware is messy; sanitize the inputs."
    *   **Responsibility**: `/bridge` directory.

---

## 🛠️ Development Workflow (Local-First)
We prioritize **Localhost execution** over containerization during the "Incubation" phase.

### Phase 1: The Vertical Slice (Current)
1.  **Start Broker**: `docker compose up -d mosquitto` (Ports: `1883`, `9001`)
2.  **Start UI**: `cd ui && npm run dev` (Localhost: `5173`)
3.  **Start Mock Sensor**: `go run ./bin/device-mock` (Publishes to `sensor/cpu/temp`)

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
