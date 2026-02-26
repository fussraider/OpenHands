package runtime

import (
	"context"
	"io"
	"openhands-go/server/config"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerRuntime struct {
	client       *client.Client
	containerID  string
	config       *config.Config
	hijackedResp *types.HijackedResponse
	shell        *ShellSession
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

func (r *DockerRuntime) ensureContainer(ctx context.Context) error {
	imageName := "ubuntu:latest" // Default, should come from config
	if r.config.Sandbox.Runtime != "" {
		// If runtime config has image name
	}

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
	return nil
}

// dockerConnWrapper wraps hijacked response to implement io.ReadWriteCloser using the BufReader
type dockerConnWrapper struct {
	resp *types.HijackedResponse
}

func (w *dockerConnWrapper) Read(p []byte) (int, error) {
	return w.resp.Reader.Read(p)
}

func (w *dockerConnWrapper) Write(p []byte) (int, error) {
	return w.resp.Conn.Write(p)
}

func (w *dockerConnWrapper) Close() error {
	w.resp.Close()
	return nil
}

func (r *DockerRuntime) startDockerShell(ctx context.Context) error {
	if r.shell != nil {
		return nil
	}

	if err := r.ensureContainer(ctx); err != nil {
		return err
	}

	// Start bash in container
	execConfig := types.ExecConfig{
		Cmd:          []string{"bash", "--noprofile", "--norc"},
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		Tty:          true,
		Env:          []string{"TERM=xterm"}, // Ensure TTY behaves correctly
	}

	execIDResp, err := r.client.ContainerExecCreate(ctx, r.containerID, execConfig)
	if err != nil {
		return err
	}

	resp, err := r.client.ContainerExecAttach(ctx, execIDResp.ID, types.ExecStartCheck{
		Tty: true,
	})
	if err != nil {
		return err
	}

	r.hijackedResp = &resp
	r.shell = NewShellSession(&dockerConnWrapper{resp: &resp})

	return nil
}

func (r *DockerRuntime) Start(ctx context.Context, command string, args ...string) error {
	// Legacy Start: launch one-off command
	if err := r.ensureContainer(ctx); err != nil {
		return err
	}

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

	resp, err := r.client.ContainerExecAttach(ctx, execIDResp.ID, types.ExecStartCheck{
		Tty: true,
	})
	if err != nil {
		return err
	}

	r.hijackedResp = &resp
	return nil
}

func (r *DockerRuntime) Execute(ctx context.Context, command string, args ...string) (string, int, error) {
	// Ensure persistent shell
	if r.shell == nil {
		if err := r.startDockerShell(ctx); err != nil {
			return "", -1, err
		}
	}

	// Check if command is bash -c which is common from ActionService
	cmdStr := command
	if command == "bash" && len(args) >= 2 && args[0] == "-c" {
		cmdStr = args[1]
	} else if len(args) > 0 {
		// Attempt to reconstruct command string
		cmdStr = command + " " + strings.Join(args, " ")
	}

	return r.shell.Execute(ctx, cmdStr)
}

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

func (r *DockerRuntime) GetCwd(ctx context.Context) (string, error) {
	if r.shell != nil {
		return r.shell.GetCwd(), nil
	}
	return "/workspace", nil
}

func (r *DockerRuntime) Close() error {
	var firstErr error
	if r.shell != nil {
		if err := r.shell.Close(); err != nil {
			firstErr = err
		}
	} else if r.hijackedResp != nil {
		r.hijackedResp.Close()
	}
	return firstErr
}
