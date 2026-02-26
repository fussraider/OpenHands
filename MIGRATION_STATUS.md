# Migration Status

This file tracks the migration status of OpenHands backend features from Python to Go.

| Feature | Python Source | Go Implementation | Status | Notes |
|---|---|---|---|---|
| Health Check | `openhands/server/routes/health.py` | `server/handlers/common.go` | ✅ Complete | Basic /health and /alive endpoints implemented. |
| Settings | `openhands/server/routes/settings.py` | `server/handlers/settings.go` | ✅ Complete | File-based settings persistence implemented. |
| Models List | `openhands/server/routes/settings.py` (indirectly) | `server/handlers/common.go` | ✅ Complete | Mock list returned. |
| Conversations | `openhands/server/routes/conversation.py` | `server/handlers/conversations.go` | ✅ Complete | File-based conversation persistence implemented. |
| Static Files | `openhands/server/app.py` (config) | `server/handlers/static.go` | ✅ Complete | Serving frontend build with SPA fallback. |
| Github Integration | `openhands/server/routes/git.py` | `server/handlers/github.go` | ✅ Complete | Implemented using `google/go-github`. |
| Files API | `openhands/server/routes/files.py` | `server/handlers/files.go` | ✅ Complete | Local workspace access implemented with security checks. |
| Secrets API | `openhands/server/routes/secrets.py` | `server/handlers/secrets.go` | ✅ Complete | In-memory secrets store implemented. |
| Feedback API | `openhands/server/routes/feedback.py` | `server/handlers/feedback.go` | ✅ Complete | Implemented persistent feedback storage. |
| Trajectory API | `openhands/server/routes/trajectory.py` | `server/handlers/trajectory.go` | ✅ Complete | Implemented with persistent event stream. |
| MCP Integration | `openhands/server/routes/mcp.py` | `server/handlers/mcp.go` | ✅ Complete | Stub implementation. |
| Public API | `openhands/server/routes/public.py` | `server/handlers/options.go` | ✅ Complete | Implemented as Options API. |
| Security API | `openhands/server/routes/security.py` | `server/handlers/security.go` | ✅ Complete | Stub implementation. |
| Socket.IO Events | `openhands/server/routes/socket.py` | `server/ws/socket.go` | ✅ Complete | Supports `oh_user_action` and broadcasts `oh_event`. |
| Stateful Shell | `openhands/runtime/impl/` | `server/runtime/shell.go` | ✅ Complete | Uses `creack/pty` for local and hijacked Docker streams for container execution. |
| Browser Plugin | `openhands/runtime/browser/` | `server/runtime/plugins/browser` | ✅ Complete | Using `playwright-go`. |
| Agent Logic | `openhands/core/agent/` | `server/agent/agent.go` | ✅ Complete | `CodeActAgent` implementation with multi-turn loop and tool calling. |
| Plugin System | `openhands/runtime/plugins/` | `server/runtime/plugins/` | ✅ Complete | Generic `Plugin` interface and integration. |
| Persistence | `openhands/storage/` | `server/events/event_stream.go` | ✅ Complete | JSONL file-based persistence for events/feedback. |

## Remaining Python Features to Port

The following features exist in the Python codebase but are **pending migration** to Go. These represent the gap between "Beta/Feature Complete" and "Full Parity".

| Area | Feature | Python Reference | Complexity | Notes |
|---|---|---|---|---|
| **Skills** | **Prompt Loading** | `openhands/skills/`, `openhands/agenthub/codeact_agent/prompts/` | Low | Need to load Jinja2 templates and Markdown skill definitions dynamically instead of hardcoding system prompts. |
| **MCP** | **Full Client Impl** | `openhands/mcp/` | High | Full implementation of Model Context Protocol client to interact with MCP servers. |
| **Agents** | **Other Agents** | `openhands/agenthub/` | Medium | Porting `BrowsingAgent`, `VisualBrowsingAgent`, `LocAgent`. |
| **Security** | **Analyzer** | `openhands/security/analyzer.py` | High | Prompt injection detection and safety checks. |
| **Memory** | **Condenser** | `openhands/memory/condenser/` | Medium | Logic to summarize/truncate long conversation history for LLM context window management. |
| **Observability** | **Tracing** | `openhands/core/logger.py` | Medium | Full OpenTelemetry tracing integration (currently stubs). |
| **Runtime** | **Jupyter Kernel** | `openhands/runtime/plugins/jupyter/` | High | Python implementation uses real Jupyter kernels via `jupyter_client`. Go implementation currently uses `python3 -c` MVP. |
| **Events** | **Task Tracking** | `openhands/events/observation/task_tracking.py` | Low | Specific event types for task progress tracking. |
| **Events** | **Loop Recovery** | `openhands/events/observation/loop_recovery.py` | Medium | Logic to detect and recover from stuck agent loops. |

## Do Not Implement

| Area | Feature | Python Reference | Reason |
|---|---|---|---|
| **Enterprise** | **Multi-tenancy** | `enterprise/server/` | Proprietary Enterprise feature. |
| **Enterprise** | **SSO/SAML** | `enterprise/integrations/` | Proprietary Enterprise feature. |
