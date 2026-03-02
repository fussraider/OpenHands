# Go Backend for OpenHands

This is the production-ready implementation of the OpenHands backend in Go, designed to replace the legacy Python backend.
It provides significantly better performance, smaller footprint, easier deployment, and stronger type safety.

## Status

**Current State:** ✅ Release Candidate (v1.0 Ready) - **100% Migration Complete**

The Go backend has reached full functional parity with the Open-Source Python implementation for all core OpenHands features.
It supports stateful shell sessions (Local and Docker runtimes), a robust plugin system (including browser capabilities via Playwright), advanced agent delegation, GitHub integration, security analysis, context condensation, and observability (OpenTelemetry).

**License Note:** This implementation strictly follows the MIT-licensed open-source codebase. Features located in the legacy Python `enterprise/` directory (under PolyForm Free Trial License) such as Multi-Tenancy and SSO are **out of scope** and intentionally omitted.

### Implemented Features

- **Frontend Serving**: Serves the React frontend build with SPA fallback routing.
- **Real-time Engine**: Full Socket.IO real-time event streaming (`oh_event`, `oh_user_action`) compatible with the React frontend.
- **Persistence & Storage**: JSONL/JSON file-based storage for `conversations.json`, `settings.json`, and trajectory event streams.
- **Security & Secrets**: Secure AES-256 encrypted at rest secrets management (`server/handlers/secrets.go`). Security Analyzer blocks risky commands (`rm -rf /`).
- **Runtime Management**:
    - **LocalRuntime**: Executes commands locally using persistent bash sessions (via `creack/pty`). Tracks `cwd` and environment variables per session.
    - **DockerRuntime**: Executes commands securely inside an isolated Docker container using the native Docker API.
    - Full file copying interfaces between host and container.
- **Agent Architecture (CodeAct)**:
    - Multi-turn Loop, Context Condensation (`TokenCondenser`), and `LoopDetector` for infinite loop recovery.
    - Advanced Tool Calling (via `tmc/langchaingo`).
    - Full parity with all core agents: `CodeActAgent`, `BrowsingAgent`, `VisualBrowsingAgent`, `ReadOnlyAgent`, and `DummyAgent`.
    - Dynamic Plugin System (`BrowserPlugin`, `JupyterPlugin`).

### API Endpoints

The following REST API endpoints are fully implemented and compatible with the v1 Frontend:

| Method | Endpoint | Description | Status |
|---|---|---|---|
| GET | `/health`, `/alive` | Server health checks | ✅ Full |
| GET | `/api/options/models` | List available LLM models | ✅ Full (Dynamic) |
| GET | `/api/options/agents` | List available agent types | ✅ Full (Dynamic Registry) |
| GET | `/api/conversations` | List conversations | ✅ Full |
| POST | `/api/conversations` | Create new conversation | ✅ Full |
| GET | `/api/conversations/{id}` | Get conversation details | ✅ Full |
| POST | `/api/conversations/{id}/action` | Execute action (legacy/rest) | ✅ Full |
| GET | `/api/conversations/{id}/list-files` | List workspace files | ✅ Full |
| GET | `/api/conversations/{id}/select-file` | Read file content | ✅ Full |
| GET | `/api/settings` | Get user settings | ✅ Full (Config merged) |
| POST | `/api/settings` | Update user settings | ✅ Full |
| GET | `/api/secrets` | List secrets | ✅ Full |
| POST | `/api/secrets` | Add secret | ✅ Full (AES-256 Encrypted) |
| GET | `/api/github/repositories` | List GitHub repos | ✅ Full |
| GET | `/api/conversations/{id}/trajectory` | Get session history | ✅ Full |
| POST | `/api/conversations/{id}/feedback` | Submit feedback | ✅ Full |
| POST/GET | `/api/v1/sandboxes/*` | Core App Server V1 Sandbox management | ✅ Full |
| POST/GET | `/api/v1/conversation/*/events/*` | Core App Server V1 Event querying | ✅ Full |
| GET | `/mcp/` | Model Context Protocol Server | ✅ Full |

---

## How to Setup and Run

### Prerequisites

- Go 1.23+
- Node.js 22+ (for building the frontend)
- Docker v28+ (required if using DockerRuntime)

### Option 1: Run Locally (Dev Mode)

1.  **Build Frontend:**
    First, compile the React UI.
    ```bash
    cd frontend
    npm install
    npm run build
    cd ..
    ```

2.  **Configure Environment (Optional but Recommended):**
    ```bash
    # Set persistent storage path for Conversations and Settings (defaults to current dir)
    export FILE_STORE_PATH=/tmp/openhands-data

    # Set an encryption key (32 bytes base64 encoded) for Secure Secrets
    export OPENHANDS_SECRETS_KEY=$(openssl rand -base64 32)
    ```

3.  **Run Go Backend:**
    ```bash
    go run cmd/server/main.go
    ```

4.  **Access UI:** Open http://localhost:3000 in your browser.

### Option 2: Run with Docker Compose (Production Ready)

This will build the complete image (Frontend + Go Backend) and run it securely.

```bash
docker-compose up --build
```
*Note: Ensure `docker-compose.yml` points to the correct Dockerfile path in the `containers/` directory.*

---

## Configuration

The server dynamically configures itself by deep-merging settings from:
1. Global config `~/.openhands/config.toml`
2. Local config `./config.toml`
3. Environment Variables (Highest Precedence)

### Key Environment Variables

| Env Var | Description | Default |
|---|---|---|
| `OPENHANDS_HOST` | Server bind address | `127.0.0.1` (`0.0.0.0` in Docker) |
| `OPENHANDS_PORT` | Server port | `3000` |
| `FILE_STORE_PATH`| Directory for JSON DB files & settings | Current directory |
| `LLM_MODEL` | Preferred LLM Model name | `gpt-4` |
| `LLM_API_KEY` | Your LLM API Key | - |
| `LLM_BASE_URL` | Custom LLM Base URL | - |
| `SANDBOX_RUNTIME`| Isolation execution mode: `local` or `docker` | `local` |
| `OPENHANDS_SECRETS_KEY` | AES-256 32-byte Base64 key for encrypting secrets | Empty (Unencrypted mode) |
| `LOG_FORMAT` | `text` or `json` for structured slog output | `text` |

---

## Testing & Observability

- **Unit & Integration Tests (Go):**
    ```bash
    # Run targeted package tests due to extensive OS-level mocking required
    go test -v ./server/agent/...
    go test -v ./server/handlers/...
    ```

- **Frontend E2E Verification:**
    ```bash
    # Ensure backend is running on port 3000
    cd frontend
    npx playwright test tests/go_backend.spec.ts
    ```

- **OpenTelemetry Tracing:**
    Tracing spans are automatically emitted by the `Agent` Loop and `LLM.CompleteWithTools` actions for local debugging or external export via OTLP.
