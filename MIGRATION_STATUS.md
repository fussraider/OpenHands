# Migration Status

This file tracks the migration status of OpenHands backend features from Python to Go.

| Feature | Python Source | Go Implementation | Status | Notes |
|---|---|---|---|---|
| Health Check | `openhands/server/routes/health.py` | `server/handlers/common.go` | ✅ Complete | Basic /health and /alive endpoints implemented. |
| Settings | `openhands/server/routes/settings.py` | `server/handlers/settings.go` | 🚧 In Progress | Initial API structure implemented. |
| Models List | `openhands/server/routes/settings.py` (indirectly) | `server/handlers/common.go` | 🚧 In Progress | Mock list returned. |
| Conversations | `openhands/server/routes/conversation.py` | `server/handlers/mock.go` | ❌ Pending | Mock list returned. |
| Static Files | `openhands/server/app.py` (config) | `server/handlers/static.go` | ✅ Complete | Serving frontend build with SPA fallback. |
| Github Integration | `openhands/server/routes/git.py` | `server/handlers/mock.go` | ❌ Pending | Mock endpoint implemented. |
| Files API | `openhands/server/routes/files.py` | `server/handlers/files.go` | ❌ Pending | |
| Feedback API | `openhands/server/routes/feedback.py` | `server/handlers/feedback.go` | ❌ Pending | |
| MCP Integration | `openhands/server/routes/mcp.py` | `server/handlers/mcp.go` | ❌ Pending | |
| Public API | `openhands/server/routes/public.py` | `server/handlers/public.go` | ❌ Pending | |
| Secrets API | `openhands/server/routes/secrets.py` | `server/handlers/secrets.go` | ❌ Pending | |
| Security API | `openhands/server/routes/security.py` | `server/handlers/security.go` | ❌ Pending | |
| Trajectory API | `openhands/server/routes/trajectory.py` | `server/handlers/trajectory.go` | ❌ Pending | |
