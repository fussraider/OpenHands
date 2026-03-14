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
		mcp.WithBoolean("draft", mcp.Description("Whether PR opened is a draft")),
	)

	createMRTool := mcp.NewTool("create_mr",
		mcp.WithDescription("Open an MR in GitLab"),
		mcp.WithString("id", mcp.Required(), mcp.Description("GitLab repository ID or URL-encoded path")),
		mcp.WithString("source_branch", mcp.Required(), mcp.Description("Source branch on repo")),
		mcp.WithString("target_branch", mcp.Required(), mcp.Description("Target branch on repo")),
		mcp.WithString("title", mcp.Required(), mcp.Description("MR Title")),
		mcp.WithString("description", mcp.Description("MR description")),
	)

	createBitbucketPRTool := mcp.NewTool("create_bitbucket_pr",
		mcp.WithDescription("Open a PR in Bitbucket"),
		mcp.WithString("repo_name", mcp.Required(), mcp.Description("Bitbucket repository (workspace/repo_slug)")),
		mcp.WithString("source_branch", mcp.Required(), mcp.Description("Source branch on repo")),
		mcp.WithString("target_branch", mcp.Required(), mcp.Description("Target branch on repo")),
		mcp.WithString("title", mcp.Required(), mcp.Description("PR Title")),
		mcp.WithString("description", mcp.Description("PR description")),
	)

	createAzurePRTool := mcp.NewTool("create_azure_devops_pr",
		mcp.WithDescription("Open a PR in Azure DevOps"),
		mcp.WithString("repo_name", mcp.Required(), mcp.Description("Azure DevOps repository")),
		mcp.WithString("source_branch", mcp.Required(), mcp.Description("Source branch on repo")),
		mcp.WithString("target_branch", mcp.Required(), mcp.Description("Target branch on repo")),
		mcp.WithString("title", mcp.Required(), mcp.Description("PR Title")),
		mcp.WithString("description", mcp.Description("PR description")),
	)

	mcpSrv.AddTool(createPRTool, handleCreatePR)
	mcpSrv.AddTool(createMRTool, handleCreateMR)
	mcpSrv.AddTool(createBitbucketPRTool, handleCreateBitbucketPR)
	mcpSrv.AddTool(createAzurePRTool, handleCreateAzurePR)

	sseServer = server.NewSSEServer(mcpSrv)
}

func getGithubToken() string {
	if token, ok := GetSecret("git_provider_github"); ok {
		return token
	}
	return ""
}

func getGitlabToken() string {
	if token, ok := GetSecret("git_provider_gitlab"); ok {
		return token
	}
	return ""
}

func getBitbucketToken() string {
	if token, ok := GetSecret("git_provider_bitbucket"); ok {
		return token
	}
	return ""
}

func getAzureToken() string {
	if token, ok := GetSecret("git_provider_azure"); ok {
		return token
	}
	return ""
}

func handleCreatePR(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}

	repoName, _ := args["repo_name"].(string)
	sourceBranch, _ := args["source_branch"].(string)
	targetBranch, _ := args["target_branch"].(string)
	title, _ := args["title"].(string)
	body, _ := args["body"].(string)
	draft, _ := args["draft"].(bool)
	labels := []string{}

	if GitService != nil {
		token := getGithubToken()
		url, err := GitService.CreatePR(ctx, token, repoName, sourceBranch, targetBranch, title, body, draft, labels)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create PR: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Successfully created PR: %s", url)), nil
	}

	return mcp.NewToolResultText("GitService not configured"), nil
}

func handleCreateMR(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}

	id, _ := args["id"].(string)
	sourceBranch, _ := args["source_branch"].(string)
	targetBranch, _ := args["target_branch"].(string)
	title, _ := args["title"].(string)
	description, _ := args["description"].(string)
	labels := []string{}

	if GitService != nil {
		token := getGitlabToken()
		url, err := GitService.CreateMR(ctx, token, id, sourceBranch, targetBranch, title, description, labels)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create MR: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Successfully created MR: %s", url)), nil
	}

	return mcp.NewToolResultText("GitService not configured"), nil
}

func handleCreateBitbucketPR(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}

	repoName, _ := args["repo_name"].(string)
	sourceBranch, _ := args["source_branch"].(string)
	targetBranch, _ := args["target_branch"].(string)
	title, _ := args["title"].(string)
	description, _ := args["description"].(string)
	labels := []string{}

	if GitService != nil {
		token := getBitbucketToken()
		url, err := GitService.CreatePR(ctx, token, repoName, sourceBranch, targetBranch, title, description, false, labels)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create PR: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Successfully created PR: %s", url)), nil
	}

	return mcp.NewToolResultText("GitService not configured"), nil
}

func handleCreateAzurePR(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args, ok := req.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("invalid arguments"), nil
	}

	repoName, _ := args["repo_name"].(string)
	sourceBranch, _ := args["source_branch"].(string)
	targetBranch, _ := args["target_branch"].(string)
	title, _ := args["title"].(string)
	description, _ := args["description"].(string)
	labels := []string{}

	if GitService != nil {
		token := getAzureToken()
		url, err := GitService.CreatePR(ctx, token, repoName, sourceBranch, targetBranch, title, description, false, labels)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create PR: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Successfully created PR: %s", url)), nil
	}

	return mcp.NewToolResultText("GitService not configured"), nil
}

func MCPSSEHandler(w http.ResponseWriter, r *http.Request) {
	sseServer.ServeHTTP(w, r)
}
