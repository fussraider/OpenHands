ARG OPENHANDS_BUILD_VERSION=dev
FROM node:22-bookworm-slim AS frontend-builder
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend ./
RUN npm run build

FROM golang:1.24-bookworm AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o openhands-server cmd/server/main.go

FROM debian:bookworm-slim AS openhands-app
ARG OPENHANDS_BUILD_VERSION

WORKDIR /app

ENV RUN_AS_OPENHANDS=true
ENV OPENHANDS_USER_ID=42420
ENV SANDBOX_LOCAL_RUNTIME_URL=http://host.docker.internal
ENV WORKSPACE_BASE=/opt/workspace_base
ENV OPENHANDS_BUILD_VERSION=$OPENHANDS_BUILD_VERSION
ENV SANDBOX_USER_ID=0
ENV FILE_STORE=local
ENV FILE_STORE_PATH=/.openhands
ENV INIT_GIT_IN_EMPTY_WORKSPACE=1

RUN mkdir -p $FILE_STORE_PATH
RUN mkdir -p $WORKSPACE_BASE

# Install runtime dependencies
# We include Python and Node.js to support code execution in LocalRuntime
RUN apt-get update -y && \
    apt-get install -y curl ssh sudo git python3 python3-pip python3-venv nodejs npm docker.io && \
    rm -rf /var/lib/apt/lists/*

# User Setup
RUN sed -i 's/^UID_MIN.*/UID_MIN 499/' /etc/login.defs && \
    sed -i 's/^UID_MAX.*/UID_MAX 1000000/' /etc/login.defs && \
    groupadd --gid $OPENHANDS_USER_ID openhands && \
    useradd -l -m -u $OPENHANDS_USER_ID --gid $OPENHANDS_USER_ID -s /bin/bash openhands && \
    usermod -aG openhands openhands && \
    usermod -aG sudo openhands && \
    echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers && \
    chown -R openhands:openhands /app && chmod -R 770 /app && \
    chown -R openhands:openhands $WORKSPACE_BASE && chmod -R 770 $WORKSPACE_BASE

USER openhands

# Copy frontend build
COPY --chown=openhands:openhands --from=frontend-builder /app/build ./frontend/build
# Copy backend binary
COPY --chown=openhands:openhands --from=backend-builder /app/openhands-server /app/openhands-server
# Copy entrypoint
COPY --chown=openhands:openhands --chmod=770 ./containers/app/entrypoint.sh /app/entrypoint.sh

ENV PATH="/app:$PATH"

USER root
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/openhands-server"]
