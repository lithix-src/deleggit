# 🤖 Deleggit Agent Swarm Assignments
> **Phase 1.5: Platform Dashboard**

## 1. UIPrime (Frontend Specialist)
**Objective**: Transform the Single-Page MVP into a Multi-Page Dashboard.

### Active Tasks
*   [ ] **Sidebar Navigation**: Implement a collapsible sidebar with links to `Dashboard`, `Workflows`, `Hardware`.
*   [ ] **Dashboard Layout**: Create a CSS Grid layout for widgets.
*   [ ] **Widget: Hardware Monitor**: Upgrade `SensorGrid` to handle CPU, RAM, and Network Latency.
*   [ ] **Widget: Active Agents**: Create a "Terminal-style" live log viewer for `agent/+/log`.
*   [ ] **Widget: Repository Events**: Create a "Notification Feed" for `repo/+/issue`.

## 2. BridgeBuilder (Hardware/Mock Specialist)
**Objective**: Simulate a complex, multi-source environment.

### Active Tasks
*   [ ] **Upgrade `device-mock`**:
    *   Refactor main loop to run multiple concurrent "emitters".
    *   **Emitter 1**: `sensor/cpu/temp` (Sine wave 40-80C).
    *   **Emitter 2**: `agent/trend-scout/log` (Random text logs: "Scanning...", "Found pattern...", "Sleeping").
    *   **Emitter 3**: `repo/lithix/issue` (Rare event: "New Issue #123 created").

## 3. SystemCore (Architect/Backend)
**Objective**: Enforce Event Consistency.

### Active Tasks
*   [ ] **Event Schema Definition**: Ensure `CloudEvents` structure is respected by both `device-mock` and `ui/src/lib/event-bus.ts`.
*   [ ] **QoS Verification**: Ensure Log messages (Agents) are delivered even if bursty.
