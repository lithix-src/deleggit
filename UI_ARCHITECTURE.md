# Catalyst UI Architecture (Dual-Mode)

To separate "Platform Operations" from "Agent Interaction", we are adopting a Dual-Mode Architecture.

## 1. Route Structure

| Route | Mode | Description | Key Components |
| :--- | :--- | :--- | :--- |
| `/` | **Mission Control** | The "Product" interface. Focused on intent and outcome. | `MissionChat`, `VibeVisualizer`, `ProjectStatus` |
| `/admin` | **Admin Console** | The "Ops" interface. Focused on telemetry and state. | `SensorGrid`, `ContainerGrid`, `RepoFeed`, `SwarmLogs` |

## 2. Layout Strategy

### A. Admin Layout (`/admin`)
*   **Style**: Dense, Data-Heavy, Sidebar Navigation.
*   **Theme**: Cyber-Industrial (Dark Slate, Monospace font).
*   **Components**: Existing Dashboard widgets.

### B. Mission Layout (`/`)
*   **Style**: Focused, Minimal, Center-Stage Chat.
*   **Theme**: "Vibe" Aesthetic (Glassmorphism, Ambient, Fluid Motion).
*   **Components**:
    *   **Chat Console**: Central input for talking to `Liaison`.
    *   **Brain Map**: Visualizer of `pgvector` recall (Nodes lighting up).
    *   **Activity Ticker**: Subtle stream of high-level events.

## 3. Technology Stack
*   **Routing**: `react-router-dom` (New Dependency).
*   **State**: `useSensorStore` (Zustand) shared between both views.
*   **Styling**: Tailwind CSS (Shared tokens, distinct layouts).

## 4. Migration Plan
1.  **Refactor**: Move `App.tsx` logic into `src/pages/AdminPage.tsx`.
2.  **Install**: `npm install react-router-dom`.
3.  **Router**: Make `App.tsx` the Router provider.
4.  **Implement**: Create `src/pages/MissionPage.tsx` (Skeleton).
