# Go Backend for OpenHands

This is a proof-of-concept implementation of the OpenHands backend in Go.

## Status

This backend implements the API structure, serves the frontend static files, and includes core agent logic.
It implements:
- **API Endpoints**: Settings, Conversations, Files, GitHub, Secrets.
- **Runtime Management**: Supports Local (PTY) and Docker runtimes.
- **Agent Logic**: Autonomous agent loop with LLM integration (using `langchaingo`).
- **Persistence**: File-based storage for settings and conversations.

Note: While functional, this is a migration work-in-progress. The agent logic is a simplified implementation compared to the Python SDK.

## Features Implemented

- **Static File Serving**: Serves the React frontend with SPA fallback.
- **API Endpoints**:
    - `/health`: Health check.
    - `/api/options/models`: Returns a mock list of models.
    - `/api/conversations`: Returns an empty list of conversations.
    - `/api/settings`: Returns default settings.
    - `/api/github/repositories`: Returns empty list.

## How to Run

1.  Build the frontend:
    ```bash
    cd frontend
    npm install
    npm run build
    cd ..
    ```

2.  Run the Go server:
    ```bash
    go run cmd/server/main.go
    ```

3.  Open http://localhost:3000 in your browser.

## Testing

- **Go Tests**:
    ```bash
    go test -v
    ```

- **Frontend Integration Test**:
    ```bash
    cd frontend
    npx playwright test tests/go_backend.spec.ts --config=playwright.go.config.ts
    ```
