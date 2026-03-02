# Backend Migration Workflows Report

## Context
The project has undergone a complete migration of its backend from Python to Go. Several GitHub Action workflows and dependencies needed to be adjusted to align with the new Go-based backend and infrastructure.
**Note:** Due to legacy requirements and tooling, the `openhands/` python directory remains in the repository and tests/linting for it have been preserved, while parallel Go workflows have been added.

## What was Modified

1. **`lint.yml`**
   - Retained the existing Python linting steps for the `openhands` python code.
   - Added a new `lint-go` job utilizing the `golangci/golangci-lint-action@v4` action to lint the Go source code with the `1.23` compiler.

2. **`lint-fix.yml`**
   - Retained the Python lint auto-fixing job.
   - Implemented a `lint-fix-go` job that performs `go fmt ./...` and runs `golangci-lint run --fix` to enforce formatting and fix lint warnings, and commits any resulting changes back to the PR.

3. **`backend-tests.yml` (New File)**
   - Created a new workflow file for testing the Go backend application running `go test ./server/... ./cmd/...`.
   - Uses standard Go `-coverprofile=coverage.out` artifact uploads.

4. **`ghcr-build.yml`**
   - Retained the Python dynamic Dockerfile generation script `openhands.runtime.utils.runtime_build` for building the tool sandbox images, as it is still relevant to the runtime environment.
   - Inserted Go `1.23` setup configuration so that the application image build steps that now compile Go binaries are successful.

5. **`openhands-resolver.yml`**
   - Updated the resolution and PR generation steps to compile and execute the Go `cmd/resolver/main.go` instead of Python `openhands.resolver.resolve_issue` scripts.
   - Adjusted dependency checks and environment setups to prepare a Go environment (`actions/setup-go@v5`) rather than `actions/setup-python`.

6. **`e2e-tests.yml`**
   - Added the Go initialization (`actions/setup-go@v5`) needed to compile and run the backend prior to E2E testing.
   - Changed cleanup steps from `pkill -f "python -m openhands.server"` to `pkill -f "server_bin"`.

7. **`dependabot.yml`**
   - Added the `gomod` ecosystem configuration to scan the Go project for dependency and security updates. (Pip is preserved since Poetry dependencies are still used by pipelines).

8. **`check-package-versions.yml`**
   - Updated the step checking `pyproject.toml` to also include a bash script check against the new `go.mod` file, ensuring no pseudo-versions or direct commits are accidentally pinned in Go dependencies.

## What was Deleted (and Why)
- Nothing was entirely deleted, as legacy Python code needs to be maintained alongside the Go backend for environment generation, tools, and scripts. All existing Python pipelines (`py-tests.yml`, etc.) were preserved while expanding the CI to support Go.
