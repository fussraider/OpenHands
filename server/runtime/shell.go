package runtime

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	cmdOutputPS1Begin = "###PS1JSON###"
	cmdOutputPS1End   = "###PS1END###"
)

var ps1Regex = regexp.MustCompile(fmt.Sprintf(`(?s)%s(.*?)%s`, regexp.QuoteMeta(cmdOutputPS1Begin), regexp.QuoteMeta(cmdOutputPS1End)))
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;?]*[a-zA-Z]`)

type ShellSession struct {
	rw          io.ReadWriteCloser
	mu          sync.Mutex
	closed      bool
	initialized bool

	// Buffer to store output
	outputBuffer bytes.Buffer

	// Current working dir cache
	workDir string
}

type CmdMetadata struct {
	ExitCode      int    `json:"exit_code,string"`
	PID           int    `json:"pid,string"`
	Username      string `json:"username"`
	Hostname      string `json:"hostname"`
	WorkingDir    string `json:"working_dir"`
	PyInterpreter string `json:"py_interpreter_path"`
}

func NewShellSession(rw io.ReadWriteCloser) *ShellSession {
	return &ShellSession{
		rw: rw,
	}
}

func (s *ShellSession) Initialize(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.initialized {
		return nil
	}

	// Disable echo to clean up output
	if _, err := s.rw.Write([]byte("stty -echo\n")); err != nil {
		return err
	}

	// Set prompt
	ps1Cmd := s.getPS1Command()
	if _, err := s.rw.Write([]byte(ps1Cmd + "\n")); err != nil {
		return err
	}

	// Wait for prompt
	// Use a short timeout for initialization
	initCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, _, err := s.readUntilPrompt(initCtx)
	if err != nil {
		return err
	}

	// Clear buffer to ensure clean slate
	s.outputBuffer.Reset()
	s.initialized = true
	return nil
}

func (s *ShellSession) getPS1Command() string {
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

func (s *ShellSession) readUntilPrompt(ctx context.Context) (string, int, error) {
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
		// Since s.rw might not support SetReadDeadline (generic io.ReadWriteCloser),
		// we rely on the underlying implementation or use a goroutine for reading if needed.
		// However, for this implementation, we assume s.rw.Read is blocking but interruptible via Close?
		// Or we need a way to timeout.
		// If s.rw is a file (PTY) or net.Conn (Docker), it supports deadlines or we can poll.
		// Since we abstracted it to io.ReadWriteCloser, we can't cast to SetReadDeadline.
		// We'll trust the caller to manage the ReadWriteCloser or use a wrapper.
		// But wait, `readUntilPrompt` needs to be responsive to ctx.Done().

		// If we do a blocking Read on s.rw, we can't respect ctx.Done() unless s.rw supports it.
		// Ideally we read in a goroutine.

		type readResult struct {
			n   int
			err error
		}
		readCh := make(chan readResult, 1)
		go func() {
			n, err := s.rw.Read(buf)
			readCh <- readResult{n, err}
		}()

		select {
		case res := <-readCh:
			if res.err != nil {
				if res.err == io.EOF {
					return "", -1, fmt.Errorf("shell exited unexpectedly")
				}
				// If it's a timeout error (from underlying implementation), we iterate.
				if os.IsTimeout(res.err) {
					continue
				}
				return "", -1, res.err
			}
			if res.n > 0 {
				s.outputBuffer.Write(buf[:res.n])
			}
		case <-ctx.Done():
			// We can't cancel the read easily if it's blocked.
			// This leaks the goroutine if Read blocks forever.
			// But for PTY/Socket, Read usually returns eventually or we Close session on cancel?
			return "", -1, ctx.Err()
		}
	}
}

func (s *ShellSession) Execute(ctx context.Context, command string) (string, int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.initialized {
		if err := s.Initialize(ctx); err != nil {
			return "", -1, err
		}
	}

	// Write command
	if _, err := s.rw.Write([]byte(command + "\n")); err != nil {
		return "", -1, err
	}

	// Read output
	rawOutput, exitCode, err := s.readUntilPrompt(ctx)
	if err != nil {
		return "", -1, err
	}

	// Clean output
	output := strings.TrimSpace(rawOutput)
	trimmedCmd := strings.TrimSpace(command)
	if strings.HasPrefix(output, trimmedCmd) {
		output = strings.TrimPrefix(output, trimmedCmd)
	}
	output = ansiRegex.ReplaceAllString(output, "")
	output = strings.TrimSpace(output)

	return output, exitCode, nil
}

func (s *ShellSession) GetCwd() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.workDir == "" {
		return "/workspace" // Default fallback
	}
	return s.workDir
}

func (s *ShellSession) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	return s.rw.Close()
}
