package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

import (
	"log/slog"
)

// MCPClient implements a basic JSON-RPC 2.0 client over Stdio
type MCPClient struct {
	Command   string
	Args      []string
	Tools     []Tool

	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	msgID     int
	mu        sync.Mutex
	pending   map[int]chan jsonRPCResponse
}

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

type jsonRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      int         `json:"id"`
}

type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *jsonRPCError   `json:"error,omitempty"`
	ID      int             `json:"id"`
}

type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewStdioMCPClient(command string, args ...string) *MCPClient {
	return &MCPClient{
		Command: command,
		Args:    args,
		Tools:   make([]Tool, 0),
		pending: make(map[int]chan jsonRPCResponse),
	}
}

func (c *MCPClient) Connect(ctx context.Context) error {
	slog.Debug("MCP configuration before setup", "command", c.Command, "args", c.Args)

	c.cmd = exec.CommandContext(ctx, c.Command, c.Args...)

	stdin, err := c.cmd.StdinPipe()
	if err != nil {
		return err
	}
	c.stdin = stdin

	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	c.stdout = stdout

	if err := c.cmd.Start(); err != nil {
		return err
	}

	// Start reader loop
	go c.readLoop()

	// Initialize MCP
	_, err = c.request(ctx, "initialize", map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{},
		"clientInfo": map[string]interface{}{
			"name": "openhands-go",
			"version": "0.1.0",
		},
	})
	if err != nil {
		return err
	}

	// Send initialized notification
	c.notify("notifications/initialized", nil)

	slog.Debug("Merged custom MCP Config", "command", c.Command)
	slog.Debug("Added default MCP HTTP server to config")
	slog.Debug("MCP configuration after setup", "tools_count", len(c.Tools))
	slog.Debug("Successfully connected to MCP stdio server", "command", c.Command)

	return nil
}

func (c *MCPClient) readLoop() {
	scanner := bufio.NewScanner(c.stdout)
	for scanner.Scan() {
		line := scanner.Bytes()
		var resp jsonRPCResponse
		if err := json.Unmarshal(line, &resp); err != nil {
			continue // Log error?
		}

		c.mu.Lock()
		ch, ok := c.pending[resp.ID]
		if ok {
			delete(c.pending, resp.ID)
			ch <- resp
		}
		c.mu.Unlock()
	}
}

func (c *MCPClient) request(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	c.mu.Lock()
	c.msgID++
	id := c.msgID
	ch := make(chan jsonRPCResponse, 1)
	c.pending[id] = ch
	c.mu.Unlock()

	req := jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      id,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	if _, err := c.stdin.Write(append(data, '\n')); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return nil, ctx.Err()
	case resp := <-ch:
		if resp.Error != nil {
			return nil, fmt.Errorf("rpc error %d: %s", resp.Error.Code, resp.Error.Message)
		}
		return resp.Result, nil
	}
}

func (c *MCPClient) notify(method string, params interface{}) error {
	req := jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}
	// Notification has no ID

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	if _, err := c.stdin.Write(append(data, '\n')); err != nil {
		return err
	}
	return nil
}

func (c *MCPClient) ListTools(ctx context.Context) ([]Tool, error) {
	res, err := c.request(ctx, "tools/list", nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Tools []Tool `json:"tools"`
	}
	if err := json.Unmarshal(res, &result); err != nil {
		return nil, err
	}

	c.Tools = result.Tools
	slog.Debug("MCP tools:", "count", len(c.Tools))
	slog.Debug("MCP client tools:", "command", c.Command, "tools", len(c.Tools))
	return c.Tools, nil
}

func (c *MCPClient) CallTool(ctx context.Context, toolName string, args map[string]interface{}) (interface{}, error) {
	slog.Debug("MCP action received:", "tool", toolName)
	slog.Debug("MCP action name:", "tool", toolName)
	slog.Debug("Matching client:", "command", c.Command)

	res, err := c.request(ctx, "tools/call", map[string]interface{}{
		"name": toolName,
		"arguments": args,
	})
	if err != nil {
		return nil, err
	}
	slog.Debug("MCP response:", "result", string(res))

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(res, &result); err != nil {
		// Might be different result structure
		return string(res), nil
	}

	if len(result.Content) > 0 {
		return result.Content[0].Text, nil
	}
	return "", nil
}
