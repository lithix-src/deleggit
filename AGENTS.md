# Catalyst Engineering Team (Automated)

> **Protocol**: `~/.gemini/docs/message-protocol.md`
> **Standards**: `docs/STANDARDS.md` (Global Engineering Standards)
> **Session**: `orch:catalyst_startup`

## Automated Functional Roles

### `spec:frontend:01` (Service: "Interface")
**Role**: Frontend Engineering
**Focus**: User Experience, Dashboard Visualization, Client Connectivity.
**Responsibility**: Maintain the `catalyst-ui` codebase and ensure operational visibility.

### `spec:backend:01` (Service: "Orchestrator")
**Role**: Core Systems Engineering
**Focus**: Event Bus, Data Persistence, Mission Routing.
**Responsibility**: Maintain `catalyst-core` and the Hexagonal Architecture adapters.

### `spec:infra:01` (Service: "Infrastructure")
**Role**: DevOps & Site Reliability
**Focus**: Kubernetes (Kind), Docker, CI/CD Pipelines.
**Responsibility**: Infrastructure-as-Code (IaC) and Environment Provisioning.

### `spec:sim:01` (Service: "Simulation")
**Role**: Synthetic Data Generation
**Focus**: Hardware Emulation, Load Testing, Chaos Engineering.
**Responsibility**: `device-mock` and generating realistic telemetry for development.

### `verif:qa:01` (Service: "Compliance")
**Role**: Quality Assurance & Audit
**Focus**: Integration Testing, Log Verification, Security Compliance.
**Responsibility**: Validating system integrity and test coverage.

## 2. Engineering Standards (The "Rigid Core")
All agents MUST adhere to the **Catalyst SDK** pattern:
1.  **Imports**: Use `github.com/point-unknown/catalyst/pkg/...` for Logging, Env, and CloudEvents.
2.  **No Stdlib Log**: Do not use `log.Println`. Use `logger.New(serviceName)`.
3.  **No Stdlib Env**: Do not use `os.Getenv`. Use `env.Get()`.
4.  **Verification**: After every code change, run `make sdk-check` to verify SDK compliance.
