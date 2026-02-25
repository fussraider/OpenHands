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

func (r *DockerRuntime) Start(ctx context.Context, command string, args ...string) error {
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
	if err := r.ensureContainer(ctx); err != nil {
		return "", -1, err
	}

	fullCmd := append([]string{command}, args...)
	execConfig := types.ExecConfig{
		Cmd:          fullCmd,
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  false,
		Tty:          true,
	}

	execIDResp, err := r.client.ContainerExecCreate(ctx, r.containerID, execConfig)
	if err != nil {
		return "", -1, err
	}

	resp, err := r.client.ContainerExecAttach(ctx, execIDResp.ID, types.ExecStartCheck{
		Tty: true,
	})
	if err != nil {
		return "", -1, err
	}
	defer resp.Close()

	output, err := io.ReadAll(resp.Reader)
	if err != nil {
		return "", -1, err
	}

	inspectResp, err := r.client.ContainerExecInspect(ctx, execIDResp.ID)
	if err != nil {
		return string(output), 0, nil
	}

	return string(output), inspectResp.ExitCode, nil
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

func (r *DockerRuntime) Close() error {
	if r.hijackedResp != nil {
		r.hijackedResp.Close()
	}
	return nil
}
