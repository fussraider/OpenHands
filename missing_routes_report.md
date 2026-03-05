# Missing Routes and Handlers in Go Backend Migration

This report details the missing routes and handlers in the Go backend compared to the existing Python backend implementation. The routes have been cross-checked, and priorities have been assigned based on their likelihood of being used by the React frontend.

## 1. High Priority (Likely blocking UI functionality)

| Python Route | HTTP Method | Python File Reference | Notes/Comment |
| --- | --- | --- | --- |
| ~~`/upload-files`~~ | POST | `openhands/server/routes/files.py:332` | **[PORTED]** Handles multiple file uploads. Included in Go backend as `POST /api/conversations/{id}/upload-files`. |
| ~~`/conversations/{id}`~~ | PATCH | `openhands/server/routes/manage_conversations.py:1142` | **[PORTED]** Renaming/updating conversation metadata (e.g. title changes). Included in Go backend. |
| `/conversations/{id}/start` | POST | `openhands/server/routes/manage_conversations.py:779` | Starts a conversation loop for an existing conversation ID. (Go has `/api/conversations` POST, but the distinct `/start` on an ID might be called). |
| `/conversations/{id}/stop` | POST | `openhands/server/routes/manage_conversations.py:850` | Pauses/Stops an active conversation loop. |
| `/message` | POST | `openhands/server/routes/conversation.py` | Standalone endpoint to send a message to a session without WebSockets. High priority if fallback polling is used. |

## 2. Medium Priority (Features that might throw soft errors or limit context)

| Python Route | HTTP Method | Python File Reference | Notes/Comment |
| --- | --- | --- | --- |
| `/suggested-tasks` | GET | `openhands/server/routes/git.py:221` | Provides suggested tasks/issues from repositories. Missing this breaks the auto-suggestion modal on the frontend UI. |
| `/conversations/{id}/microagents` | GET | `openhands/server/routes/manage_conversations.py` | Specific conversation microagents payload. Go only has the global user microagents endpoint (`/api/user/repository/.../microagents`). |
| `/conversations/{id}/remember-prompt` | GET | `openhands/server/routes/manage_conversations.py:671` | Fetches a prompt generated from the memory condenser/microagents. |
| `/vscode-url` | GET | `openhands/server/routes/conversation.py:150` | Returns the URL for the embedded VSCode workspace. Will break the embedded IDE button/iframe. |

## 3. Low Priority (Deprecated, redundant, or internal)

| Python Route | HTTP Method | Python File Reference | Notes/Comment |
| --- | --- | --- | --- |
| `/web-hosts` | GET | `openhands/server/routes/conversation.py:185` | Port mapping details. May be legacy. |
| `/conversations/{id}/exp-config` | POST | `openhands/server/routes/manage_conversations.py` | Experimental configuration endpoint. Usually non-critical. |
| `/events` | GET, POST | `openhands/server/routes/conversation.py` | The legacy global event fetching. Mostly replaced by `/api/conversations/{id}/trajectory` and `/api/v1/conversation/{id}/events`. |
| `/config` | GET | `openhands/server/routes/conversation.py` / `public.py` | Redundant. Go handles configuration via `/api/settings` and `/api/v1/web-client/config`. |
| `/server_info`, `/ready` | GET | `openhands/server/routes/health.py` | Go implements `/health` and `/alive`, which handles typical liveness probes adequately. |

## Existing/Implemented Handlers (For Reference)

The following groups of routes are confirmed to be successfully ported and functional in Go:
* `GET /api/options/*` (Models, Agents, Security Analyzers)
* `GET/POST /api/conversations` (Core conversation creation and fetching)
* `GET/DELETE /api/conversations/{id}`
* `GET/POST/PUT/DELETE /api/secrets`
* `GET/POST /api/settings`
* `GET /api/v1/web-client/config` (and other v1 sandbox mocks)
* `GET/POST /api/user/*` (Git installations, repos, branches, search)
* `GET /api/conversations/{id}/list-files`, `select-file`
* `GET /api/conversations/{id}/trajectory`
* `POST /api/conversations/{id}/action`, `submit-feedback`
