# Go Backend for OpenHands

This is a proof-of-concept implementation of the OpenHands backend in Go, designed to replace the existing Python backend.
It aims to provide better performance, easier deployment, and stronger type safety.

## Status

**Current State:** 🚧 Beta / Feature Complete

This backend implements the core API structure, serves the frontend static files, and includes a fully functional autonomous agent loop.
It supports stateful shell sessions, agent delegation, plugin system, and GitHub integration.

**License Note:** This implementation strictly follows the MIT-licensed open-source codebase. Features located in the `enterprise/` directory (under PolyForm Free Trial License) are **out of scope** and not implemented here.

### Implemented Features

- **Frontend Serving**: Serves the React frontend build with SPA fallback.
- **Socket.IO**: Full real-time event streaming (`oh_event`, `oh_user_action`) compatible with the frontend.
- **Persistence**: File-based storage for:
    - Conversations (`conversations.json`)
    - Settings (`settings.json`)
    - Secrets (In-memory for now)
- **Runtime Management**:
    - **LocalRuntime**: Executes commands locally using persistent bash session (via `creack/pty`). Supports stateful execution (cwd, env vars).
    - **DockerRuntime**: Executes commands inside a Docker container using the Docker API.
- **Agent Logic**:
    - Basic "Loop" that fetches events and queries an LLM.
    - Integration with `tmc/langchaingo` for LLM support.
    - Supports `RUN` (execute command) and `MSG` (chat) actions.

### API Endpoints

The following REST API endpoints are implemented:

| Method | Endpoint | Description | Status |
|---|---|---|---|
| GET | `/health`, `/alive` | Server health checks | ✅ Full |
| GET | `/api/options/models` | List available LLM models | ✅ Mock |
| GET | `/api/options/agents` | List available agent types | ✅ Mock |
| GET | `/api/conversations` | List conversations | ✅ Full |
| POST | `/api/conversations` | Create new conversation | ✅ Full |
| GET | `/api/conversations/{id}` | Get conversation details | ✅ Full |
| POST | `/api/conversations/{id}/action` | Execute action (legacy/rest) | ✅ Full |
| GET | `/api/conversations/{id}/list-files` | List workspace files | ✅ Full |
| GET | `/api/conversations/{id}/select-file` | Read file content | ✅ Full |
| GET | `/api/settings` | Get user settings | ✅ Full |
| POST | `/api/settings` | Update user settings | ✅ Full |
| GET | `/api/secrets` | List secrets | ✅ Full |
| POST | `/api/secrets` | Add secret | ✅ Full |
| GET | `/api/github/repositories` | List GitHub repos | ✅ Full |
| GET | `/api/conversations/{id}/trajectory` | Get session history | ✅ Full |
| POST | `/api/conversations/{id}/feedback` | Submit feedback | 🚧 Stub |
| GET | `/mcp/` | Model Context Protocol | 🚧 Stub |

## How to Run

### Prerequisites

- Go 1.23+
- Node.js 22+ (for building frontend)
- Docker (optional, for DockerRuntime)

### Option 1: Run Locally (Dev Mode)

1.  **Build Frontend:**
    ```bash
    cd frontend
    npm install
    npm run build
    cd ..
    ```

2.  **Run Backend:**
    ```bash
    # Set persistent storage path (optional, defaults to current dir)
    export FILE_STORE_PATH=/tmp/openhands-data

    # Run server
    go run cmd/server/main.go
    ```

3.  **Access:** Open http://localhost:3000

### Option 2: Run with Docker Compose

This will build the complete image (Frontend + Go Backend) and run it.

```bash
docker-compose up --build
```
*Note: Ensure `docker-compose.yml` points to the correct Dockerfile or update the build context if needed.*

## Configuration

The server is configured via `config.toml` or Environment Variables.
Env vars take precedence.

| Env Var | Description | Default |
|---|---|---|
| `OPENHANDS_HOST` | Server bind address | `127.0.0.1` (or `0.0.0.0` in Docker) |
| `OPENHANDS_PORT` | Server port | `3000` |
| `FILE_STORE_PATH`| Directory for JSON DB files | Current directory |
| `LLM_MODEL` | LLM Model name | `gpt-4` |
| `LLM_API_KEY` | LLM API Key | - |
| `LLM_BASE_URL` | Custom LLM Base URL | - |
| `SANDBOX_RUNTIME`| `local` or `docker` | `local` |

## Testing

- **Unit Tests (Go):**
    ```bash
    go test ./...
    ```

- **Integration Tests (Frontend + Backend):**
    ```bash
    # Ensure backend is running on port 3000
    cd frontend
    npx playwright test tests/go_backend.spec.ts
    ```
