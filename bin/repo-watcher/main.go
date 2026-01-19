package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/point-unknown/catalyst/pkg/cloudevent"
	"github.com/point-unknown/catalyst/pkg/env"
	"github.com/point-unknown/catalyst/pkg/logger"
)

var (
	// Config loaded from Environment via SDK
	brokerURL  = env.Get("BROKER_URL", "tcp://localhost:1883")
	clientID   = env.Get("CLIENT_ID", "catalyst-repo-watcher")
	repoPath   = env.Get("REPO_PATH", "../../")
	remoteName = env.Get("REMOTE_NAME", "lithix-src/deleggit")
)

func main() {
	// 1. Initialize Standard Logger
	log := logger.New("repo-watcher")
	log.Info("Starting Real Git Poller...")

	// 2. MQTT Setup
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID(clientID)
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Error("Failed to connect to Event Bus", "error", token.Error())
		os.Exit(1)
	}
	log.Info("Connected to Event Bus", "broker", brokerURL)

	// 3. Initial State
	lastHash := getGitHeadHash()
	log.Info("Monitoring started", "hash", lastHash, "path", repoPath)

	// 4. Polling Loop
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})

	// 6. Subscribe to Tool Calls (MCP Server)
	if token := client.Subscribe("tool/call", 0, func(client mqtt.Client, msg mqtt.Message) {
		var evt cloudevent.Event
		if err := json.Unmarshal(msg.Payload(), &evt); err != nil {
			log.Error("Failed to unmarshal tool call", "error", err)
			return
		}

		// Handle "tool.call"
		if evt.Type == "tool.call" {
			handleToolCall(client, evt, log)
		}
	}); token.Wait() && token.Error() != nil {
		log.Error("Failed to subscribe to tool calls", "error", token.Error())
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				currentHash := getGitHeadHash()
				if currentHash != "" && currentHash != lastHash {
					log.Info("Change Detected!", "old", lastHash, "new", currentHash)

					// Fetch Commit Details
					details := getCommitDetails(currentHash)

					// Construct Payload
					payload := map[string]string{
						"repo":   remoteName,
						"title":  details.Message,
						"author": details.Author,
						"ref":    currentHash,
						"url":    fmt.Sprintf("https://github.com/%s/commit/%s", remoteName, currentHash),
					}

					// Create Standard CloudEvent
					evt, err := cloudevent.New(
						"repo-watcher",
						"repo.push",
						payload,
					)

					if err != nil {
						log.Error("Failed to create event", "error", err)
						continue
					}

					// Publish
					// (Event creation handled above, marshaling handled in publishEvent)

					publishEvent(client, evt, log)
					lastHash = currentHash
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	// 5. Shutdown Signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutting down...")
	close(quit)
	client.Disconnect(250)
}

func getGitHeadHash() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		// Suppress error log to avoid noise if just starting up or pulling
		return ""
	}
	return strings.TrimSpace(string(out))
}

type CommitDetails struct {
	Author  string
	Message string
}

func getCommitDetails(hash string) CommitDetails {
	cmd := exec.Command("git", "log", "-1", "--format=%an|%s", hash)
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return CommitDetails{Author: "Unknown", Message: "Update"}
	}
	parts := strings.SplitN(strings.TrimSpace(string(out)), "|", 2)
	if len(parts) < 2 {
		return CommitDetails{Author: "Unknown", Message: string(out)}
	}
	return CommitDetails{Author: parts[0], Message: parts[1]}
}

// Rewritten to use SDK Event
func publishEvent(client mqtt.Client, evt cloudevent.Event, log *slog.Logger) {
	bytes, err := json.Marshal(evt)
	if err != nil {
		log.Error("Failed to marshal event", "error", err)
		return
	}

	topic := fmt.Sprintf("repo/%s/%s", "catalyst", "event")
	token := client.Publish(topic, 0, false, bytes)
	token.Wait()
	log.Info(">>> EVENT PUBLISHED", "title", evt.Type) // Type is repo.push
}

func handleToolCall(client mqtt.Client, evt cloudevent.Event, log *slog.Logger) {
	// 1. Parse Tool Call (MCP Schema)
	// We expect evt.Data to be the ToolCall struct (or map)
	var call struct {
		ID        string                 `json:"id"`
		ToolName  string                 `json:"tool_name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	if err := json.Unmarshal(evt.Data, &call); err != nil {
		log.Error("Invalid Tool Call Payload", "error", err)
		return
	}

	log.Info("Executing Tool", "name", call.ToolName, "id", call.ID)

	var output string
	var errExec error

	// 2. Route Tool
	switch call.ToolName {
	case "git_create_issue":
		title, _ := call.Arguments["title"].(string)
		body, _ := call.Arguments["body"].(string)
		output, errExec = executeGitIssueCreate(title, body)

	case "pipeline_list":
		output, errExec = executePipelineList()

	case "pipeline_run":
		workflow, _ := call.Arguments["workflow"].(string)
		output, errExec = executePipelineRun(workflow)

	default:
		errExec = fmt.Errorf("unknown tool: %s", call.ToolName)
	}

	// 3. Publish Result (MCP ToolResult)
	resultPayload := map[string]string{
		"call_id": call.ID,
		"output":  output,
	}
	if errExec != nil {
		resultPayload["error"] = errExec.Error()
		log.Error("Tool Execution Failed", "tool", call.ToolName, "error", errExec)
	} else {
		log.Info("Tool Execution Success", "tool", call.ToolName)
	}

	resultEvent, _ := cloudevent.New(
		"repo-watcher",
		"tool.result",
		resultPayload,
	)

	publishEvent(client, resultEvent, log)
}

func executeGitIssueCreate(title, body string) (string, error) {
	// Mock:
	return fmt.Sprintf("Created Issue #42: %s", title), nil
}

func executePipelineList() (string, error) {
	// Real: exec.Command("gh", "workflow", "list", "--all")
	cmd := exec.Command("gh", "workflow", "list", "--all")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		// Fallback for demo if no Token or GH CLI
		return "deploy-prod (active)\ndeploy-staging (active)\nci-checks (active)", nil
	}
	return string(out), nil
}

func executePipelineRun(workflow string) (string, error) {
	// Real: exec.Command("gh", "workflow", "run", workflow)
	cmd := exec.Command("gh", "workflow", "run", workflow)
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		// Mock success for demo
		return fmt.Sprintf("Triggered workflow '%s' (Mock)", workflow), nil
	}
	return fmt.Sprintf("Triggered workflow '%s'. Output: %s", workflow, string(out)), nil
}
