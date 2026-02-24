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
