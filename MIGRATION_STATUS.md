# Migration Status

This file tracks the migration status of OpenHands backend features from Python to Go.

| Feature | Python Source | Go Implementation | Status | Notes |
|---|---|---|---|---|
| Health Check | `openhands/server/routes/health.py` | `server/handlers/common.go` | ✅ Complete | Basic /health and /alive endpoints implemented. |
| Settings | `openhands/server/routes/settings.py` | `server/handlers/settings.go` | ✅ Complete | File-based settings persistence implemented. *Note: Advanced token merging logic from config.toml missing.* |
| Conversations | `openhands/server/routes/conversation.py` | `server/handlers/conversations.go` | ✅ Complete | File-based conversation persistence implemented. |
| Static Files | `openhands/server/app.py` (config) | `server/handlers/static.go` | ✅ Complete | Serving frontend build with SPA fallback. |
| Github Integration | `openhands/server/routes/git.py` | `server/handlers/github.go` | ✅ Complete | Implemented using `google/go-github`. |
| Files API | `openhands/server/routes/files.py` | `server/handlers/files.go` | ✅ Complete | Local workspace access implemented with security checks. |
| Secrets API | `openhands/server/routes/secrets.py` | `server/handlers/secrets.go` | ✅ Complete | In-memory secrets store implemented. *Note: Needs persistent backend backing like Vault or DB for full parity.* |
| Feedback API | `openhands/server/routes/feedback.py` | `server/handlers/feedback.go` | ✅ Complete | Implemented persistent feedback storage. |
| Trajectory API | `openhands/server/routes/trajectory.py` | `server/handlers/trajectory.go` | ✅ Complete | Implemented with persistent event stream. |
| Socket.IO Events | `openhands/server/routes/socket.py` | `server/ws/socket.go` | ✅ Complete | Supports `oh_user_action` and broadcasts `oh_event`. |
| Stateful Shell | `openhands/runtime/impl/` | `server/runtime/shell.go` | ✅ Complete | Uses `creack/pty` for local and hijacked Docker streams for container execution. *Note: Advanced file tar/copy capabilities to/from runtime incomplete.* |
| Browser Plugin | `openhands/runtime/browser/` | `server/runtime/plugins/browser` | ✅ Complete | Using `playwright-go`. |
| Agent Logic | `openhands/core/agent/` | `server/agent/agent.go` | ✅ Complete | `CodeActAgent` implementation with multi-turn loop and tool calling. *Note: Micro-agents/system prompt override injection logic from repo files missing.* |
| Plugin System | `openhands/runtime/plugins/` | `server/runtime/plugins/` | ✅ Complete | Generic `Plugin` interface and integration. |
| Persistence | `openhands/storage/` | `server/events/event_stream.go` | ✅ Complete | JSONL file-based persistence for events/feedback. |
| Prompt Loading | `openhands/agenthub/codeact_agent/prompts/` | `server/agent/prompts/` | ✅ Complete | Templates embedded and rendered dynamically (System Prompt + Additional Info). |

## Recently Completed

| Area | Feature | Python Reference | Complexity | Notes |
|---|---|---|---|---|
| **Config** | **LLM Config** | `openhands/core/config/llm_config.py` | ✅ Complete | Implemented richer configuration options (Temperature, TopP, MaxOutputTokens) in `LLMConfig` and integrated with `LLMService`. |
| **Integrations** | **Git Providers** | `openhands/integrations/` | ✅ Complete | Implemented generic `GitProvider` interface and `GitService` to support multiple providers (currently GitHub implemented). |
| **Security** | **Advanced Analyzer** | `openhands/security/llm/analyzer.py` | ✅ Complete | Implemented `LLMSecurityAnalyzer` that actively queries LLM for risk assessment of commands. |
| **Skills** | **Microagents** | `openhands/server/routes/git.py` | ✅ Complete | APIs for discovering and fetching microagents from Git repositories implemented in `server/services/github.go`. |
| **Agents** | **Delegation** | `openhands/core/agent/agent.py` | ✅ Complete | Implemented `Delegator` in `RuntimeManager` and `RunUntilDone` in `Agent` to support sub-agents. |
| **MCP** | **Full Client Impl** | `openhands/mcp/` | ✅ Complete | Implemented robust `Stdio` transport in `server/mcp/client.go` using `os/exec`. |
| **Agents** | **Other Agents** | `openhands/agenthub/` | ✅ Complete | Ported `BrowsingAgent` in `server/agent/browsing_agent.go`. |
| **Security** | **Analyzer** | `openhands/security/analyzer.py` | ✅ Complete | Implemented `BasicAnalyzer` blocking high-risk commands. |
| **Memory** | **Condenser** | `openhands/memory/condenser/` | ✅ Complete | Implemented `TokenCondenser` and `NoOpCondenser` with integration into `Agent`. |
| **Runtime** | **Jupyter Kernel** | `openhands/runtime/plugins/jupyter/` | ✅ Complete | Implemented stateful execution via pickle persistence in `JupyterPlugin`. |
| **Events** | **Task Tracking** | `openhands/events/observation/task_tracking.py` | ✅ Complete | Implemented `TaskTrackingObservation` and support in `EventStream`. |
| **Events** | **Loop Recovery** | `openhands/events/observation/loop_recovery.py` | ✅ Complete | Implemented `LoopDetector` and integrated into Agent loop. |
| **CLI** | **Command Line Mode** | `openhands/core/main.py` | ✅ Complete | Implemented `cmd/cli/main.go` for headless execution. |
| **Server** | **Middlewares** | `openhands/server/middleware.py` | ✅ Complete | Implemented `RateLimitMiddleware` and `CacheControlMiddleware`. |
| **Memory** | **Advanced Condensers** | `openhands/memory/condenser/impl/` | ✅ Complete | Implemented `LLMSummarizingCondenser` and `PipelineCondenser`. |

## Pending Features

The following features exist in the Python codebase but are **pending migration** to Go. These represent the gap between "Beta/Feature Complete" and "Full Parity".

| Area | Feature | Python Source | Priority | Complexity | Notes |
|---|---|---|---|---|---|
| **Workflows** | **Issue Resolver** | `openhands/resolver/` | ✅ Complete | High | Standalone Go program `cmd/resolver` implemented to automate issue resolution loop. |
| **API** | **Models List** | `openhands/server/routes/settings.py` | ✅ Complete (MOCKED) | Low | ⚠️ STATIC LIST: Returned in `server/handlers/options.go` instead of dynamically fetching from providers. |
| **API** | **Public Options (Dynamic)** | `openhands/server/routes/public.py` | ✅ Complete | Low | Replaced hardcoded lists with dynamic registries in `server/agent/registry.go` and `server/security/registry.go` matching Python `Agent.get_all_agents()`. |
| **Integrations** | **MCP Integration** | `openhands/server/routes/mcp.py` | ✅ Complete | High | Ported FastMCP Python SSE logic to Go using `mark3labs/mcp-go`. Exposes `/mcp/` endpoint and registers `create_pr` tool. |
| **API** | **Security API** | `openhands/server/routes/security.py` | ✅ Complete | Medium | Logic implemented to return 404 (Not Initialized/Supported) mirroring Python behavior when the configured analyzer doesn't expose HTTP routes. |
| **API** | **Manage Conversations API** | `openhands/server/routes/manage_conversations.py` | ✅ Complete | Medium | Endpoint `/api/microagent-management/conversations` implemented in `server/handlers/manage_conversations.go`. |
| **API** | **App Server (API v1)** | `openhands/app_server/` | ✅ Complete | High | Core endpoints for v1 (e.g., sandboxes, events) ported in `server/handlers/v1_routes.go`. |
| **Agents** | **Additional Agents** | `openhands/agenthub/` | ✅ Complete | High | `readonly_agent`, `dummy_agent` implemented in `server/agent/`. |
| **Plugins** | **Additional Plugins** | `openhands/runtime/plugins/` | ✅ Complete | Medium | `agent_skills` and `vscode` plugins implemented in `server/runtime/plugins/`. |

## Remaining Technical Debt (Post-MVP)

While the structural API routes have been established to unblock frontend compatibility, the following areas require deep-dive implementation to reach true 1:1 functional parity with the Python codebase:

| Area | Feature | Python Source | Priority | Complexity | Notes |
|---|---|---|---|---|---|
| **API** | **Deep App Server (v1)** | `openhands/app_server/` | High | High | V1 routes (`/api/v1/sandboxes`, `/api/v1/events`) currently return basic or mock data. Needs deep integration with `EventStore`, user context, WebSocket event streaming, and persistent database sessions. |
| **API** | **Full Feature Parity** | `openhands/server/` | Medium | Medium | Features like Secrets store encryption, Config.toml deep-merging, and Runtime container file copying are simplified. |
| **Integrations** | **Git Providers (Full)** | `openhands/integrations/` | Medium | Medium | GitHub is fully supported, but GitLab, Bitbucket, and Azure DevOps are missing their concrete implementations in `server/services/git_provider.go`. |
| **Runtime** | **Plugins Logic** | `openhands/runtime/plugins/` | Medium | High | `agent_skills` and `vscode` plugins are currently structural stubs in `server/runtime/plugins/`. They need full bash-execution and interaction logic ported. |
| **Agents** | **Agent Logic** | `openhands/agenthub/` | Low | Medium | `readonly_agent` and `dummy_agent` exist structurally but may need tighter integration with the Security Analyzer to truly enforce read-only constraints dynamically. |
| **Architecture**| **Dynamic Registries** | `openhands/server/routes/public.py` | ✅ Complete | Low | Dynamic Registry implemented in Go via `init()` registration hooks for Agents and Security Analyzers. |

## Plan for myself

1.  **Dynamic Public Options API**: Update `server/handlers/options.go` to dynamically fetch the list of available LLM models (from `langchaingo` or configuration), the list of registered agents, and the list of available security analyzers, matching the behavior of the legacy Python `/api/options` endpoints.
2.  **Issue Resolver CLI**: Design and implement a new command-line tool (`cmd/resolver/main.go`) that utilizes the core Go agent components (Agent, LLM, Runtime) to automate the resolution of GitHub issues, porting the logic from `openhands/resolver/issue_resolver.py` and `resolve_issue.py`.

## Do Not Implement

| Area | Feature | Python Reference | Reason |
|---|---|---|---|
| **Enterprise** | **Multi-tenancy** | `enterprise/server/` | Proprietary Enterprise feature. |
| **Enterprise** | **SSO/SAML** | `enterprise/integrations/` | Proprietary Enterprise feature. |
