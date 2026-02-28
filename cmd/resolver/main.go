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
)

func main() {
	issueNumber := flag.Int("issue", 0, "GitHub Issue Number to resolve")
	repo := flag.String("repo", "", "GitHub Repository (owner/repo)")
	token := flag.String("token", os.Getenv("GITHUB_TOKEN"), "GitHub Token")
	flag.Parse()

	logger.Init()

	if err := config.LoadConfig(); err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	if *issueNumber == 0 || *repo == "" || *token == "" {
		fmt.Println("Error: Please provide an issue number (-issue), repository (-repo), and token (-token or GITHUB_TOKEN env).")
		flag.Usage()
		os.Exit(1)
	}

	ctx := context.Background()

	// 1. Fetch Issue Details

	// We need to fetch the issue. Since GitProvider doesn't have GetIssue,
	// we will construct a specialized prompt for the agent to fetch it using bash/curl if needed,
	// or we can use the Github client directly if we expose it.
	// For this CLI, the easiest way to achieve the "Resolver" workflow without adding new GitProvider methods
	// is to instruct the agent to use the GitHub CLI (gh) or curl to read the issue and then fix it.

	task := fmt.Sprintf(`Please resolve GitHub Issue #%d in repository %s.
1. Use the 'gh' CLI or curl with the provided GITHUB_TOKEN to fetch the issue details.
2. Clone the repository and checkout a new branch.
3. Fix the issue described.
4. Commit your changes, push the branch, and open a Pull Request.
`, *issueNumber, *repo)

	sid := "resolver-" + uuid.New().String()

	// Initialize Runtime Manager
	rm := services.NewRuntimeManager()

	// Create Runtime
	_, err := rm.CreateRuntime(ctx, sid)
	if err != nil {
		slog.Error("Failed to create runtime", "error", err)
		os.Exit(1)
	}
	defer rm.StopRuntime(sid)

	// Inject GITHUB_TOKEN into the runtime environment for the agent to use
	// We do this by sending an initial command to export it.

	// Create EventStream
	es := events.NewEventStream(sid, "")

	// Set env var in runtime
	es.AddEvent(events.Event{
		ID:   uuid.New().String(),
		Type: events.EventTypeAction,
		Content: models.CmdRunAction{
			Action:  models.ActionTypeCmdRun,
			Command: fmt.Sprintf("export GITHUB_TOKEN=%s", *token),
		},
		Source: "user",
	})

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

	slog.Info("Starting Agent via Resolver CLI", "sid", sid, "issue", *issueNumber, "repo", *repo)

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
	fmt.Println("\nIssue Resolution Process Completed.")
}
