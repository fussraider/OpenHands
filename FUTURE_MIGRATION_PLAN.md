# Future Migration Plan: Python to Go

This document outlines the roadmap for addressing the remaining technical debt and achieving full feature parity in the migration of the OpenHands backend from Python to Go. It is based on the remaining tasks identified in `MIGRATION_STATUS.md`.

**Critical constraint:** All future work must respect the open-source license. No features from the `enterprise/` directory (e.g., multi-tenancy, SSO/SAML) should be ported. Complete compatibility with the existing React frontend must be maintained.

---

## 1. Dynamic Public Options API

Currently, `/api/options/models` in `server/handlers/options.go` returns a static/mock list of models.
*   **Goal:** Dynamically fetch the list of available LLM models.
*   **Implementation Steps:**
    1.  Investigate `tmc/langchaingo` capabilities or the underlying provider's API (e.g., OpenAI, Anthropic) to fetch available models dynamically.
    2.  Update `server/handlers/options.go` (`ModelsHandler`) to merge the configured model with the dynamically fetched list.
    3.  Ensure the response format remains identical to the existing frontend expectations.
    *Note: Agents and Security Analyzers are already dynamically fetched via registries, so only the models list needs to be implemented.*

## 2. Deep App Server (v1)

Currently, the V1 routes (`/api/v1/sandboxes`, `/api/v1/events`) in `server/handlers/v1_routes.go` return basic or mock data.
*   **Goal:** Replace mock data with fully functional endpoints supporting the frontend's advanced requirements, minus any enterprise-specific logic.
*   **Implementation Steps:**
    1.  **Event Store Integration:** Connect the v1 event endpoints to the `EventStore`/`EventStream` system to provide real event history and real-time streaming via websockets/SSE.
    2.  **Sandbox Management:** Implement real sandbox lifecycle management (Start, Pause, Resume, Delete) utilizing the `RuntimeManager`. Map user/session contexts to specific runtime instances.
    3.  **User Context & DB Sessions:** Integrate a persistent session/database mechanism for V1 routes (while explicitly avoiding multi-tenancy features found in the enterprise version).
    4.  Verify all V1 endpoints format responses exactly as the legacy Python `app_server` did, as the frontend depends on this structure.

## 3. Full Feature Parity: Secrets, Config, and File Copying

Several core utilities are simplified in the Go MVP.
*   **Goal:** Achieve 1:1 functional parity with the Python implementation for these utilities.
*   **Implementation Steps:**
    1.  **Secrets Store Encryption:** Update `server/handlers/secrets.go` to encrypt secrets at rest (e.g., using AES encryption with a master key derived from the environment), replacing the current plaintext in-memory map. Ensure persistence across restarts if required by parity.
    2.  **Config Deep-Merging:** Update `server/config/config.go` to support deep merging of `config.toml` options with environment variables and defaults, fully matching the complex configuration loading logic of the Python backend.
    3.  **Runtime File Copying:** Enhance `server/runtime/shell.go` (and `DockerRuntime` / `LocalRuntime`) to support advanced file tar/copy capabilities (e.g., `CopyFileToContainer`, `CopyFileFromContainer`). This is necessary for robust agent workspace manipulation.

## 4. Git Providers (Full)

GitHub is currently fully supported, but other providers are missing concrete implementations.
*   **Goal:** Support all git providers available in the Python codebase.
*   **Implementation Steps:**
    1.  Review `openhands/integrations/` in the Python source to understand the required provider abstractions.
    2.  In `server/services/git_provider.go` (or equivalent package), implement the `GitProvider` interface for:
        *   **GitLab**
        *   **Bitbucket**
        *   **Azure DevOps**
    3.  Ensure authentication, repository listing, and PR creation (if applicable) work for all implemented providers. Update secrets handling to manage tokens for these providers.

## 5. Runtime Plugins Logic

The `agent_skills` and `vscode` plugins exist as structural stubs in `server/runtime/plugins/`.
*   **Goal:** Port full bash-execution and interaction logic for these plugins.
*   **Implementation Steps:**
    1.  **Agent Skills (`openhands/runtime/plugins/agent_skills`):** Implement the logic to inject the skills library into the runtime, making standard functions (like file editing, searching) available to the agent's bash session.
    2.  **VSCode (`openhands/runtime/plugins/vscode`):** Implement the initialization and execution logic required to support VSCode integration within the runtime environment, matching Python behavior.

## 6. Agent Logic Refinement

`readonly_agent` and `dummy_agent` exist structurally but need refinement.
*   **Goal:** Ensure these agents enforce their specific constraints strictly and dynamically.
*   **Implementation Steps:**
    1.  **Readonly Agent:** Tighter integration with the Security Analyzer (`server/security/`) is needed to truly enforce read-only constraints dynamically across all tools and bash execution.
    2.  **Dummy Agent:** Review `dummy_agent` logic to guarantee it correctly mirrors the Python counterpart's minimal implementation and response formats.

## 7. Issue Resolver CLI Enhancements

The foundational structure for `cmd/resolver/main.go` exists, utilizing the core Go agent components.
*   **Goal:** Automate the resolution of GitHub issues, porting the full logic from `openhands/resolver/issue_resolver.py` and `resolve_issue.py`.
*   **Implementation Steps:**
    1.  Review the existing Go implementation against the Python source.
    2.  Implement robust branch management, PR creation (leveraging the new Git Provider implementations), and handling of resolver output formats (`resolver_output.py`, `visualize_resolver_output.py`).
    3.  Add comprehensive logging and error handling to ensure reliability in headless automation scenarios.
