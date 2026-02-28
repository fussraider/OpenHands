package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"openhands-go/server/config"
	"openhands-go/server/events"
	"openhands-go/server/logger"
	"openhands-go/server/models"
	"openhands-go/server/services"
	"os"

	"github.com/google/uuid"
	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
	"strings"
)

func main() {
	issueNumber := flag.Int("issue", 0, "GitHub Issue Number to resolve")
	repoStr := flag.String("repo", "", "GitHub Repository (owner/repo)")
	token := flag.String("token", os.Getenv("GITHUB_TOKEN"), "GitHub Token")
	flag.Parse()

	logger.Init()

	if err := config.LoadConfig(); err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	if *issueNumber == 0 || *repoStr == "" || *token == "" {
		fmt.Println("Error: Please provide an issue number (-issue), repository (-repo), and token (-token or GITHUB_TOKEN env).")
		flag.Usage()
		os.Exit(1)
	}

	parts := strings.Split(*repoStr, "/")
	if len(parts) != 2 {
		fmt.Println("Error: Repository must be in the format 'owner/repo'")
		os.Exit(1)
	}
	owner, repoName := parts[0], parts[1]

	ctx := context.Background()

	// 1. Fetch Issue Details programmatically
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	issue, _, err := client.Issues.Get(ctx, owner, repoName, *issueNumber)
	if err != nil {
		slog.Error("Failed to fetch issue from GitHub", "error", err)
		os.Exit(1)
	}

	issueTitle := issue.GetTitle()
	issueBody := issue.GetBody()

	slog.Info("Fetched issue", "title", issueTitle)

	// We pass the issue context via prompt without placing the token in the event stream history
	task := fmt.Sprintf(`Please resolve the following GitHub Issue in repository %s/%s.

ISSUE TITLE: %s
ISSUE BODY:
%s

Instructions:
1. Clone the repository locally: git clone https://github.com/%s/%s.git
2. cd into the repository and checkout a new branch 'fix-issue-%d'
3. Implement the necessary changes to fix the issue described above.
4. Verify your changes.
`, owner, repoName, issueTitle, issueBody, owner, repoName, *issueNumber)

	sid := "resolver-" + uuid.New().String()

	// Initialize Runtime Manager
	rm := services.NewRuntimeManager()

	// Create Runtime
	rt, err := rm.CreateRuntime(ctx, sid)
	if err != nil {
		slog.Error("Failed to create runtime", "error", err)
		os.Exit(1)
	}
	defer rm.StopRuntime(sid)

	// Inject GITHUB_TOKEN directly into the runtime securely (via runtime API)
	// We use rt.Execute to run it silently inside the sandbox without logging it to event stream
	// Note: LocalRuntime and DockerRuntime execute commands in bash.
	_, _, err = rt.Execute(ctx, "export GITHUB_TOKEN='" + *token + "'")
	if err != nil {
		slog.Error("Failed to inject token into runtime", "error", err)
		os.Exit(1)
	}

	// Create EventStream
	es := events.NewEventStream(sid, "")

	// Add Initial Task
	es.AddEvent(events.Event{
		ID:   uuid.New().String(),
		Type: events.EventTypeAction,
		Content: models.MessageAction{
			Action:  models.ActionTypeMessage,
			Content: task,
		},
		Source: "user",
	})

	// Add event subscriber to print events to console
	es.Subscribe(func(e events.Event) {
		if e.Source == "agent" {
			if msgAct, ok := e.Content.(models.MessageAction); ok {
				fmt.Printf("\n[AGENT]: %s\n", msgAct.Content)
			} else if runAct, ok := e.Content.(models.CmdRunAction); ok {
				fmt.Printf("\n[AGENT] Executing: %s\n", runAct.Command)
			} else if finAct, ok := e.Content.(models.AgentFinishAction); ok {
				fmt.Printf("\n[AGENT] Finished: %+v\n", finAct)
			}
		} else if e.Source == "runtime" {
			if obs, ok := e.Content.(models.CmdOutputObservation); ok {
				fmt.Printf("\n[RUNTIME Output]: %s\n", obs.Content)
			}
		}
	})

	slog.Info("Starting Agent via Resolver CLI", "sid", sid, "issue", *issueNumber, "repo", *repoStr)

	// Start Agent
	err = rm.StartAgent(ctx, sid, es)
	if err != nil {
		slog.Error("Failed to start agent", "error", err)
		os.Exit(1)
	}

	// Wait for AgentFinishAction
	done := make(chan struct{})
	es.Subscribe(func(e events.Event) {
		if e.Type == events.EventTypeAction && e.Source == "agent" {
			if _, ok := e.Content.(models.AgentFinishAction); ok {
				close(done)
			}
		}
	})

	<-done
	fmt.Println("\nAgent finished editing. Issue Resolution Process Completed.")

	// Optional: programmatic PR creation could happen here.
	// For now, the agent has made the code changes locally.
}
