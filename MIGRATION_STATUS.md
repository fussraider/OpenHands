# Migration Status

This file tracks the migration status of OpenHands backend features from Python to Go.

| Feature | Python Source | Go Implementation | Status | Notes |
|---|---|---|---|---|
| Health Check | `openhands/server/routes/health.py` | `server/handlers/common.go` | ✅ Complete | Basic /health and /alive endpoints implemented. |
| Settings | `openhands/server/routes/settings.py` | `server/handlers/settings.go` | ✅ Complete | File-based settings persistence implemented. |
| Models List | `openhands/server/routes/settings.py` (indirectly) | `server/handlers/common.go` | ✅ Complete | Mock list returned. |
| Conversations | `openhands/server/routes/conversation.py` | `server/handlers/conversations.go` | ✅ Complete | File-based conversation persistence implemented. |
| Static Files | `openhands/server/app.py` (config) | `server/handlers/static.go` | ✅ Complete | Serving frontend build with SPA fallback. |
| Github Integration | `openhands/server/routes/git.py` | `server/handlers/github.go` | 🚧 In Progress | Mock endpoint implemented. |
| Files API | `openhands/server/routes/files.py` | `server/handlers/files.go` | ✅ Complete | Local workspace access implemented with security checks. |
| Secrets API | `openhands/server/routes/secrets.py` | `server/handlers/secrets.go` | ✅ Complete | In-memory secrets store implemented. |
| Feedback API | `openhands/server/routes/feedback.py` | `server/handlers/feedback.go` | ✅ Complete | Stub implementation. |
| Trajectory API | `openhands/server/routes/trajectory.py` | `server/handlers/trajectory.go` | ✅ Complete | Stub implementation. |
| MCP Integration | `openhands/server/routes/mcp.py` | `server/handlers/mcp.go` | ✅ Complete | Stub implementation. |
| Public API | `openhands/server/routes/public.py` | `server/handlers/options.go` | ✅ Complete | Implemented as Options API. |
| Security API | `openhands/server/routes/security.py` | `server/handlers/security.go` | ✅ Complete | Stub implementation. |
| Socket.IO Events | `openhands/server/routes/socket.py` | `server/ws/socket.go` | ✅ Complete | Supports `oh_user_action` and broadcasts `oh_event`. |

## Future Improvements & Missing Features

The following features exist in the Python codebase but are not yet implemented or are incomplete in the Go backend.
**Important:** Do not implement features from `enterprise/` directory as they are under a restrictive license.

| Area | Feature | Priority | Python Reference | Notes |
|---|---|---|---|---|
| **Runtime** | **Stateful Shell Execution** | ✅ Complete | `openhands/runtime/impl/` | `LocalRuntime` now supports persistent bash sessions via `creack/pty`. `DockerRuntime` supports `Execute` interface but currently stateless (one-off exec). |
| **Runtime** | **Browser Automation** | 🟡 High | `openhands/runtime/browser/` | Integration with Playwright/Selenium for web browsing agents. |
| **Runtime** | **Plugin System** | 🟡 High | `openhands/runtime/plugins/` | Support for dynamic plugin loading (e.g., linters, formatters). |
| **Agent** | **Complex Agent Logic** | ✅ Complete | `openhands/core/agent/` | Implemented `CodeActAgent` logic with tool calling (`execute_bash`) and multi-turn reasoning using `langchaingo`. |
| **Agent** | **Delegation** | 🟡 Medium | `openhands/agenthub/` | Support for delegating tasks to sub-agents. |
| **Events** | **Rich Event Types** | 🟡 Medium | `openhands/events/` | Expand `Event` struct to support specific Action/Observation types (e.g., `CmdRunAction`, `CmdOutputObservation`) with proper validation. |
| **API** | **Real GitHub Auth** | 🟡 Medium | `openhands/server/routes/git.py` | Implement actual OAuth flow for GitHub integration. |
| **Security** | **Sandboxing** | 🔴 Critical | `openhands/runtime/` | Ensure robust isolation for `LocalRuntime` (e.g., using Firejail or similar if running outside Docker). |
| **Observability** | **Structured Logging** | 🟢 Low | `openhands/core/logger.py` | Replace standard `log` with structured logger (e.g., `slog` or `zap`). |
| **Enterprise** | **Multi-tenancy** | ⛔ Out of Scope | `enterprise/server/` | **DO NOT IMPLEMENT**. Enterprise feature. |
| **Enterprise** | **SSO/SAML** | ⛔ Out of Scope | `enterprise/integrations/` | **DO NOT IMPLEMENT**. Enterprise feature. |
