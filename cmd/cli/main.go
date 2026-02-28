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
	task := flag.String("t", "", "The task for the agent to perform")
	sessionID := flag.String("sid", "", "Session ID (optional)")
	flag.Parse()

	logger.Init()

	if err := config.LoadConfig(); err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	if *task == "" {
		fmt.Println("Error: Please provide a task using -t flag.")
		os.Exit(1)
	}

	sid := *sessionID
	if sid == "" {
		sid = "cli-" + uuid.New().String()
	}

	ctx := context.Background()

	// Initialize Runtime Manager
	rm := services.NewRuntimeManager()

	// Create Runtime
	_, err := rm.CreateRuntime(ctx, sid)
	if err != nil {
		slog.Error("Failed to create runtime", "error", err)
		os.Exit(1)
	}
	defer rm.StopRuntime(sid)

	// Create EventStream
	es := events.NewEventStream(sid, "") // In-memory for CLI unless config overrides

	// Add Initial Task
	es.AddEvent(events.Event{
		ID:   uuid.New().String(),
		Type: events.EventTypeAction,
		Content: models.MessageAction{
			Action:  models.ActionTypeMessage,
			Content: *task,
		},
		Source: "user",
	})

	// Add CLI event subscriber to print events to console
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

	slog.Info("Starting Agent via CLI", "sid", sid, "task", *task)

	// Start Agent
	err = rm.StartAgent(ctx, sid, es)
	if err != nil {
		slog.Error("Failed to start agent", "error", err)
		os.Exit(1)
	}

	// Wait for AgentFinishAction
	// We can poll the event stream or use a channel.
	// Since we are running StartAgent (which runs loop in background), we need to block main.
	// We can subscribe specifically for finish.
	done := make(chan struct{})
	es.Subscribe(func(e events.Event) {
		if e.Type == events.EventTypeAction && e.Source == "agent" {
			if _, ok := e.Content.(models.AgentFinishAction); ok {
				close(done)
			}
		}
	})

	<-done
	fmt.Println("\nTask completed.")
}
