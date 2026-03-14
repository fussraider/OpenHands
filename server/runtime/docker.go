package runtime

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"log/slog"
	"openhands-go/server/config"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
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
	imageName := "ubuntu:latest" // Default
	if r.config.Sandbox.ContainerImage != "" {
		imageName = r.config.Sandbox.ContainerImage
	}

	if r.containerID == "" {
		slog.Debug("Pulling Docker image (if needed)", "image", imageName)
		reader, err := r.client.ImagePull(ctx, imageName, image.PullOptions{})
		if err == nil {
			io.Copy(io.Discard, reader)
			reader.Close()
		} else {
			slog.Warn("Failed to pull image, will attempt to use local", "error", err)
		}

		workspaceDir := os.Getenv("WORKSPACE_BASE")
		if workspaceDir == "" {
			cwd, _ := os.Getwd()
			workspaceDir = filepath.Join(cwd, "workspace")
		}
		os.MkdirAll(workspaceDir, 0755)

		resp, err := r.client.ContainerCreate(ctx, &container.Config{
			Image:      imageName,
			Cmd:        []string{"tail", "-f", "/dev/null"}, // Keep alive
			Tty:        true,
			WorkingDir: "/workspace",
		}, &container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: workspaceDir,
					Target: "/workspace",
				},
			},
		}, nil, nil, "")
		if err != nil {
			return err
		}
		r.containerID = resp.ID

		if err := r.client.ContainerStart(ctx, r.containerID, container.StartOptions{}); err != nil {
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
	execConfig := container.ExecOptions{
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

	resp, err := r.client.ContainerExecAttach(ctx, execIDResp.ID, container.ExecAttachOptions{
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
	execConfig := container.ExecOptions{
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

	resp, err := r.client.ContainerExecAttach(ctx, execIDResp.ID, container.ExecAttachOptions{
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

	slog.Debug("Executing command:", "command", cmdStr)
	out, exitCode, err := r.shell.Execute(ctx, cmdStr)
	slog.Debug("Command finished", "exit_code", exitCode)

	return out, exitCode, err
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

func (r *DockerRuntime) CopyFileToContainer(ctx context.Context, hostPath string, containerPath string) error {
	file, err := os.Open(hostPath)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	// Create an in-memory tar buffer
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	hdr := &tar.Header{
		Name: filepath.Base(containerPath),
		Mode: 0644,
		Size: stat.Size(),
	}

	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}

	if _, err := io.Copy(tw, file); err != nil {
		return err
	}

	if err := tw.Close(); err != nil {
		return err
	}

	dir := filepath.Dir(containerPath)
	return r.client.CopyToContainer(ctx, r.containerID, dir, &buf, container.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
}

func (r *DockerRuntime) CopyFileFromContainer(ctx context.Context, containerPath string, hostPath string) error {
	reader, _, err := r.client.CopyFromContainer(ctx, r.containerID, containerPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	tr := tar.NewReader(reader)
	_, err = tr.Next() // Get first file in archive
	if err != nil {
		return err
	}

	destFile, err := os.Create(hostPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, tr)
	return err
}

func (r *DockerRuntime) GetVSCodeURL() *string {
	// Not implemented natively in basic Go Docker driver yet.
	return nil
}

func (r *DockerRuntime) GetWebHosts() map[string]interface{} {
	// Not implemented natively in basic Go Docker driver yet.
	return make(map[string]interface{})
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

	if r.containerID != "" {
		// Remove the container
		err := r.client.ContainerRemove(context.Background(), r.containerID, container.RemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		})
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
