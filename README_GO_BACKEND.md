# Go Backend for OpenHands

This is a proof-of-concept implementation of the OpenHands backend in Go.

## Status

This backend currently implements the API structure and serves the frontend static files.
It **does not** yet implement the core agent logic, runtime management, or LLM integration.
These components are complex and tightly coupled with the Python SDK.

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
