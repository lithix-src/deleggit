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
