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

## Do Not Implement

| Area | Feature | Python Reference | Reason |
|---|---|---|---|
| **Enterprise** | **Multi-tenancy** | `enterprise/server/` | Proprietary Enterprise feature. |
| **Enterprise** | **SSO/SAML** | `enterprise/integrations/` | Proprietary Enterprise feature. |
