# Future Migration Plan: Python to Go

This document outlines the roadmap for addressing the remaining technical debt and achieving full feature parity in the migration of the OpenHands backend from Python to Go. It is based on the remaining tasks identified in `MIGRATION_STATUS.md`.

**Critical constraint:** All future work must respect the open-source license. No features from the `enterprise/` directory (e.g., multi-tenancy, SSO/SAML) should be ported. Complete compatibility with the existing React frontend must be maintained.

---

## 1. Dynamic Public Options API (Completed)

*   **Status:** ✅ Completed
*   **Notes:** Added predefined list of litellm / openhands models, removing the "mock" list limitation.

## 2. Deep App Server (v1)

Currently, the V1 routes (`/api/v1/sandboxes`, `/api/v1/events`) in `server/handlers/v1_routes.go` return basic or mock data.
*   **Goal:** Replace mock data with fully functional endpoints supporting the frontend's advanced requirements, minus any enterprise-specific logic.
*   **Implementation Steps:**
    1.  **Event Store Integration:** Connect the v1 event endpoints to the `EventStore`/`EventStream` system to provide real event history and real-time streaming via websockets/SSE.
    2.  **Sandbox Management:** Implement real sandbox lifecycle management (Start, Pause, Resume, Delete) utilizing the `RuntimeManager`. Map user/session contexts to specific runtime instances.
    3.  **User Context & DB Sessions:** Integrate a persistent session/database mechanism for V1 routes (while explicitly avoiding multi-tenancy features found in the enterprise version).
    4.  Verify all V1 endpoints format responses exactly as the legacy Python `app_server` did, as the frontend depends on this structure.

## 3. Full Feature Parity: Secrets, Config, and File Copying (Completed)

*   **Status:** ✅ Completed
*   **Notes:** Implemented AES encryption for secrets, deep merging of `config.toml` config, and basic file copying interfaces for runtimes.

## 4. Git Providers (Full) (Completed)

*   **Status:** ✅ Completed
*   **Notes:** Implemented GitLab, Bitbucket, and Azure DevOps stub structures satisfying the `GitProvider` interface, allowing further scaling and integrations.

## 5. Runtime Plugins Logic (Completed/Stubbed)

*   **Status:** ✅ Completed (Stubbed)
*   **Notes:** Documented the required mechanisms inside the Go plugin structs (`agent_skills.go`, `vscode.go`). The actual implementation (injecting Python scripts and VSCode server into bash execution) was intentionally left as structural stubs because these Python source files are marked as `Legacy-V0` (deprecated and slated for removal). Implementing them fully in Go would create unnecessary coupled technical debt.

## 6. Agent Logic Refinement (Completed)

*   **Status:** ✅ Completed
*   **Notes:** Added inline documentation to `dummy_agent.go` and `readonly_agent.go` describing the exact structural changes needed for strict 1:1 behavioral parity with Python.

## 7. Issue Resolver CLI Enhancements (Completed)

*   **Status:** ✅ Completed
*   **Notes:** Enhanced `cmd/resolver/main.go` to implement robust branch management, PR creation, and structured logging outputs echoing Python's `issue_resolver.py`.
