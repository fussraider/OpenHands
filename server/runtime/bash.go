package runtime

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/creack/pty"
)

const (
	cmdOutputPS1Begin = "###PS1JSON###"
	cmdOutputPS1End   = "###PS1END###"
)

var ps1Regex = regexp.MustCompile(fmt.Sprintf(`(?s)%s(.*?)%s`, regexp.QuoteMeta(cmdOutputPS1Begin), regexp.QuoteMeta(cmdOutputPS1End)))
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;?]*[a-zA-Z]`)

type BashSession struct {
	cmd     *exec.Cmd
	pty     *os.File
	mu      sync.Mutex
	workDir string
	closed  bool

	// Buffer to store output
	outputBuffer bytes.Buffer
}

type CmdMetadata struct {
	ExitCode int    `json:"exit_code,string"`
	PID      int    `json:"pid,string"`
	Username string `json:"username"`
	Hostname string `json:"hostname"`
	WorkingDir string `json:"working_dir"`
	PyInterpreter string `json:"py_interpreter_path"`
}

func NewBashSession(workDir string) (*BashSession, error) {
	return &BashSession{
		workDir: workDir,
	}, nil
}

func (s *BashSession) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.start()
}

func (s *BashSession) start() error {
	if s.cmd != nil && s.cmd.Process != nil {
		return nil // Already running
	}

	// Ensure workDir exists
	if s.workDir != "" {
		if err := os.MkdirAll(s.workDir, 0755); err != nil {
			return err
		}
	}

	s.cmd = exec.Command("bash", "--noprofile", "--norc")
	if s.workDir != "" {
		s.cmd.Dir = s.workDir
	}
	s.cmd.Env = os.Environ()
	// Disable prompt expansion for PS1 initially to avoid noise? No.

	f, err := pty.Start(s.cmd)
	if err != nil {
		return err
	}
	s.pty = f

	// Disable echo to clean up output
	if _, err := s.pty.Write([]byte("stty -echo\n")); err != nil {
		return err
	}

	// Set prompt
	ps1Cmd := s.getPS1Command()
	if _, err := s.pty.Write([]byte(ps1Cmd + "\n")); err != nil {
		return err
	}

	// Wait for prompt
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, _, err = s.readUntilPrompt(ctx)
	// Clear buffer to ensure clean slate
	s.outputBuffer.Reset()
	return err
}

func (s *BashSession) getPS1Command() string {
	// We break the markers to avoid matching the command echo.
	// markerBegin generates ###PS1JSON###
	markerBegin := "$(printf \"%s%s\" \"###PS1\" \"JSON###\")"
	markerEnd := "$(printf \"%s%s\" \"###PS1\" \"END###\")"

	// JSON Template (minimal)
	// We use single quotes for PS1 assignment to delay expansion and avoid double-quote escaping hell.
	jsonTemplate := `{"exit_code": "$?", "working_dir": "$(pwd)"}`

	// export PS1='...'
	return fmt.Sprintf("export PS1='%s%s%s'; export PS2=\"\"", markerBegin, jsonTemplate, markerEnd)
}

func (s *BashSession) readUntilPrompt(ctx context.Context) (string, int, error) {
	buf := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			return "", -1, ctx.Err()
		default:
		}

		// Check if prompt is ALREADY in buffer
		data := s.outputBuffer.Bytes()
		loc := ps1Regex.FindIndex(data)
		if loc != nil {
			// Found
			rawOutput := string(data[:loc[0]])
			matchStr := string(data[loc[0]:loc[1]]) // metadata including markers

			// Extract metadata
			matches := ps1Regex.FindStringSubmatch(matchStr)
			exitCode := -1
			if len(matches) > 1 {
				var meta CmdMetadata
				if err := json.Unmarshal([]byte(matches[1]), &meta); err == nil {
					exitCode = meta.ExitCode
					if meta.WorkingDir != "" {
						s.workDir = meta.WorkingDir
					}
				}
			}

			// Update buffer: keep remaining
			remaining := data[loc[1]:]
			remCopy := make([]byte, len(remaining))
			copy(remCopy, remaining)
			s.outputBuffer.Reset()
			s.outputBuffer.Write(remCopy)

			return rawOutput, exitCode, nil
		}

		// Read more
		s.pty.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		n, err := s.pty.Read(buf)
		if err != nil {
			if err == io.EOF {
				return "", -1, fmt.Errorf("shell exited unexpectedly")
			}
			if !os.IsTimeout(err) {
				return "", -1, err
			}
			// Timeout -> continue loop
			continue
		}
		if n > 0 {
			s.outputBuffer.Write(buf[:n])
		}
	}
}

func (s *BashSession) Execute(ctx context.Context, command string) (string, int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cmd == nil || s.cmd.Process == nil {
		if err := s.start(); err != nil {
			return "", -1, err
		}
	}

	// Write command
	if _, err := s.pty.Write([]byte(command + "\n")); err != nil {
		return "", -1, err
	}

	// Read output
	rawOutput, exitCode, err := s.readUntilPrompt(ctx)
	if err != nil {
		return "", -1, err
	}

	// Clean output: remove the command echo if present (if stty -echo failed or slow)
	output := strings.TrimSpace(rawOutput)
	trimmedCmd := strings.TrimSpace(command)
	if strings.HasPrefix(output, trimmedCmd) {
		output = strings.TrimPrefix(output, trimmedCmd)
	}
	// Remove ANSI escape codes (e.g. bracketed paste)
	output = ansiRegex.ReplaceAllString(output, "")
	output = strings.TrimSpace(output)

	return output, exitCode, nil
}

func (s *BashSession) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	if s.pty != nil {
		s.pty.Close()
	}
	if s.cmd != nil && s.cmd.Process != nil {
		return s.cmd.Process.Kill()
	}
	return nil
}
