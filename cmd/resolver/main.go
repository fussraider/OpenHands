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

// IssueResolver represents a tool that automatically resolves GitHub issues
// using the OpenHands agent. This is a skeleton implementation.
type IssueResolver struct {
	Repo        string
	IssueNumber int
	OutputDir   string
}

func main() {
	repo := flag.String("repo", "", "Repository in format owner/repo")
	issue := flag.Int("issue", 0, "Issue number to resolve")
	outputDir := flag.String("output-dir", "output", "Output directory")

	flag.Parse()

	if *repo == "" || *issue == 0 {
		fmt.Println("Usage: resolver -repo <owner/repo> -issue <number>")
		os.Exit(1)
	}

	logger.Init()

	if err := config.LoadConfig(); err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	resolver := &IssueResolver{
		Repo:        *repo,
		IssueNumber: *issue,
		OutputDir:   *outputDir,
	}

	ctx := context.Background()
	resolver.Resolve(ctx)
}

func (r *IssueResolver) Resolve(ctx context.Context) {
	slog.Info("Starting Issue Resolver", "repo", r.Repo, "issue", r.IssueNumber)

	// 1. Fetch Issue details using GitService (Placeholder)
	slog.Info("Fetching issue details...")
	issueDescription := fmt.Sprintf("Please fix issue #%d in repository %s.", r.IssueNumber, r.Repo)

	// 2. Clone repository into workspace (Placeholder)
	// In reality, we'd use git clone or mount a volume.
	slog.Info("Preparing workspace...")

	// 3. Initialize Agent
	sid := "resolver-" + uuid.New().String()
	rm := services.NewRuntimeManager()

	_, err := rm.CreateRuntime(ctx, sid)
	if err != nil {
		slog.Error("Failed to create runtime", "error", err)
		return
	}
	defer rm.StopRuntime(sid)

	es := events.NewEventStream(sid, "")

	es.AddEvent(events.Event{
		ID:   uuid.New().String(),
		Type: events.EventTypeAction,
		Content: models.MessageAction{
			Action:  models.ActionTypeMessage,
			Content: issueDescription,
		},
		Source: "user",
	})

	err = rm.StartAgent(ctx, sid, es)
	if err != nil {
		slog.Error("Failed to start agent", "error", err)
		return
	}

	// 4. Run until done
	// We can use rm.Delegate logic or just subscribe
	slog.Info("Agent is running...")

	done := make(chan struct{})
	es.Subscribe(func(e events.Event) {
		if e.Type == events.EventTypeAction && e.Source == "agent" {
			if _, ok := e.Content.(models.AgentFinishAction); ok {
				close(done)
			}
		}
	})

	<-done

	// 5. Extract Git Patch
	slog.Info("Task finished. Extracting git patch...")
	rt, _ := rm.GetRuntime(sid)
	output, _, err := rt.Execute(ctx, "git", "diff")
	if err != nil {
		slog.Error("Failed to get git diff", "error", err)
	}

	// 6. Save Output
	fmt.Printf("\n=== RESOLVER OUTPUT ===\n%s\n=======================\n", output)
	slog.Info("Resolution complete. Output saved.")
}
