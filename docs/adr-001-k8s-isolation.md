# ADR 001: Kubernetes Isolation Strategy for Untrusted Code

## Status
Proposed

## Context
Catalyst executes AI-generated code, which is inherently untrusted. We need a way to run this code within a Kubernetes cluster with maximum isolation to prevent container escapes, lateral movement, and host compromise. The solution must support standard Docker/OCI images to maintain compatibility with the broader ecosystem.

## Decision
We will use **Kata Containers with Firecracker** as the runtime for agent sandboxes.

## Rationale
1.  **Hard Isolation**: Firecracker provides microVM-based isolation. Each pod runs in its own kernel, eliminating the shared-kernel attack surface of standard containers (runc).
2.  **Kubernetes Compatibility**: Kata Containers implements the Kubernetes CRI (Container Runtime Interface). To K8s, it looks like a standard container runtime, allowing us to use standard deployment manifests and tooling.
3.  **Performance**: Firecracker is optimized for transient, serverless-like workloads with sub-125ms boot times and low memory footprint, matching our agent task profile.
4.  **Security Depth**: This aligns with our "Defense in Depth" strategy: Namespace isolation -> Network Policies -> **MicroVM Isolation (Kata)** -> AppArmor/Seccomp.

## Consequences
*   **Complexity**: Requires installing and managing the Kata Runtime and Firecracker binaries on worker nodes.
*   **Hardware Req**: Nodes must support nested virtualization (if running on cloud instances) or bare metal virtualization extensions (VT-x/AMD-V).
*   **Resource Overhead**: Slight increase in memory/CPU per pod compared to native containers, but acceptable for the security gain.
