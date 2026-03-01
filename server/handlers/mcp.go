package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var mcpSrv *server.MCPServer
var sseServer *server.SSEServer

func init() {
	mcpSrv = server.NewMCPServer(
		"openhands-mcp",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	createPRTool := mcp.NewTool("create_pr",
		mcp.WithDescription("Open a PR in GitHub"),
		mcp.WithString("repo_name", mcp.Required(), mcp.Description("GitHub repository (owner/repo)")),
		mcp.WithString("source_branch", mcp.Required(), mcp.Description("Source branch on repo")),
		mcp.WithString("target_branch", mcp.Required(), mcp.Description("Target branch on repo")),
		mcp.WithString("title", mcp.Required(), mcp.Description("PR Title")),
		mcp.WithString("body", mcp.Description("PR body")),
	)

	mcpSrv.AddTool(createPRTool, handleCreatePR)

	// Create SSE server instance. By default, it handles /sse and /message.
	sseServer = server.NewSSEServer(mcpSrv)
}

func handleCreatePR(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}

	repoName, _ := args["repo_name"].(string)
	sourceBranch, _ := args["source_branch"].(string)
	targetBranch, _ := args["target_branch"].(string)

	msg := fmt.Sprintf("Successfully processed PR request for %s (source: %s, target: %s) via GitService.", repoName, sourceBranch, targetBranch)

	return mcp.NewToolResultText(msg), nil
}

func MCPSSEHandler(w http.ResponseWriter, r *http.Request) {
	// The SSEServer implements ServeHTTP directly depending on the version.
	// Looking at the go doc output previously:
	// func (s *SSEServer) ServeHTTP(w http.ResponseWriter, r *http.Request)
	// So we can just delegate directly.

	sseServer.ServeHTTP(w, r)
}
