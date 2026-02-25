package runtime

import (
	"context"
	"io"
	"openhands-go/server/config"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerRuntime struct {
	client       *client.Client
	containerID  string
	config       *config.Config
	hijackedResp *types.HijackedResponse
}

func NewDockerRuntime(cfg *config.Config) (*DockerRuntime, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerRuntime{
		client: cli,
		config: cfg,
	}, nil
}

func (r *DockerRuntime) Start(ctx context.Context, command string, args ...string) error {
	// For DockerRuntime in this context (OpenHands), "Start" typically means
	// ensuring the sandbox container is running.
	// However, the interface implies running a command.
	// If the container isn't running, we should start it.

	imageName := "ubuntu:latest" // Default, should come from config
	if r.config.Sandbox.Runtime != "" {
		// If runtime config has image name (TODO: add to config)
	}

	// 1. Create/Start Container if not exists
	if r.containerID == "" {
		resp, err := r.client.ContainerCreate(ctx, &container.Config{
			Image: imageName,
			Cmd:   []string{"tail", "-f", "/dev/null"}, // Keep alive
			Tty:   true,
		}, nil, nil, nil, "")
		if err != nil {
			return err
		}
		r.containerID = resp.ID

		if err := r.client.ContainerStart(ctx, r.containerID, types.ContainerStartOptions{}); err != nil {
			return err
		}
	}

	// 2. Execute Command
	// The current Runtime interface assumes Start runs a command and holds a pipe.
	// This maps poorly to Docker Exec which is one-off.
	// But let's implement it as Exec.

	fullCmd := append([]string{command}, args...)
	execConfig := types.ExecConfig{
		Cmd:          fullCmd,
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		Tty:          true,
	}

	execIDResp, err := r.client.ContainerExecCreate(ctx, r.containerID, execConfig)
	if err != nil {
		return err
	}

	// Attach to exec
	resp, err := r.client.ContainerExecAttach(ctx, execIDResp.ID, types.ExecStartCheck{
		Tty: true,
	})
	if err != nil {
		return err
	}

	// We need to store resp to allow Read/Write/Close
	// But DockerRuntime struct fields for IO?
	// The Runtime interface design (single instance per command vs per session) is slightly ambiguous here.
	// If Runtime instance = Session, then Start() running a command implies the session runs ONE command?
	// Refactoring the interface might be better, but sticking to plan:
	// We'll store the hijacker response.

	// This is a simplification. Real OpenHands uses a persistent agent server inside Docker.
	// Here we just wrap `docker exec`.

	// Hack: We can't easily store the HijackedResponse in the struct if we want to support multiple commands
	// unless we create a new Runtime instance per command (which ActionService does effectively).
	// But `runtimeManager` keeps one Runtime per conversation.
	// So `Start` should probably *not* be used for individual commands if the Runtime represents the *Environment*.

	// ADJUSTMENT: The `Runtime` interface in `interface.go` has `Start(ctx, command, args)`.
	// This implies the Runtime object *is* the process.
	// If so, `RuntimeManager` should create a new Runtime object for each Action?
	// Or `RuntimeManager` manages the *Environment* (Container) and we need an `Executor` abstraction?

	// For now, let's assume `DockerRuntime` struct represents the *Active Command* inside a container.
	// And we need a factory or the `NewDockerRuntime` should assume an existing container?

	// Let's implement it such that `Start` executes the command in a (potentially shared) container.
	// We need to manage the container lifecycle separately or check if it exists.

	r.hijackedResp = &resp
	return nil
}

// Add field to struct
// hijackedResp *types.HijackedResponse
// But types.HijackedResponse is a struct, not interface.

func (r *DockerRuntime) Write(p []byte) (n int, err error) {
	if r.hijackedResp == nil {
		return 0, io.ErrClosedPipe
	}
	return r.hijackedResp.Conn.Write(p)
}

func (r *DockerRuntime) Read(p []byte) (n int, err error) {
	if r.hijackedResp == nil {
		return 0, io.ErrClosedPipe
	}
	return r.hijackedResp.Reader.Read(p)
}

func (r *DockerRuntime) Close() error {
	if r.hijackedResp != nil {
		r.hijackedResp.Close()
	}
	// Don't kill container here if we want persistence, but for "Process" semantics, we might want to kill the Exec?
	// Docker doesn't easily support killing an Exec without killing container or finding PID.
	return nil
}
