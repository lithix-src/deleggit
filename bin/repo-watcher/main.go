package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	brokerURL = "tcp://localhost:1883"
	clientID  = "catalyst-repo-watcher"
	topics    = []string{
		"repo/lithix-src/catalyst/push",
		"repo/lithix-src/catalyst/issue",
		"repo/lithix-src/catalyst/pr",
	}
)

// GitHubEvent simulates a GitHub API response
type GitHubEvent struct {
	Type   string
	Repo   string
	Title  string
	Ref    string
	Author string
}

func main() {
	log.Println("[RepoWatcher] Starting Smart Poller...")

	// 1. MQTT Setup
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID(clientID)
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("[RepoWatcher] Failed to connect: %v", token.Error())
	}
	log.Println("[RepoWatcher] Connected to Event Bus")

	// 2. State Tracking
	lastSeenSHA := "initial-sha" // In real app, load from DB

	// 3. Polling Loop
	// We poll frequently (every 5s), but only EMIT on change.
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				// SIMULATE: Check for changes
				// In a real implementation, this would be:
				// commits, err := github.ListCommits(since=lastSeenSHA)

				if shouldTriggerChange() {
					evt := generateRandomEvent()

					// Dedup check (Simulated)
					if evt.Ref == lastSeenSHA {
						continue
					}
					lastSeenSHA = evt.Ref

					publishEvent(client, evt)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	// 4. Shutdown Signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("[RepoWatcher] Shutting down...")
	close(quit)
	client.Disconnect(250)
}

// shouldTriggerChange simulates the rarity of real events.
// Returns true only 20% of the time to avoid spam.
func shouldTriggerChange() bool {
	return rand.Float32() < 0.2
}

func generateRandomEvent() GitHubEvent {
	eventTypes := []string{"push", "issue", "pr"}
	chosenType := eventTypes[rand.Intn(len(eventTypes))]

	ref := fmt.Sprintf("sha-%d", time.Now().UnixNano())

	switch chosenType {
	case "push":
		return GitHubEvent{
			Type:   "repo.push",
			Repo:   "lithix-src/catalyst",
			Title:  fmt.Sprintf("feat: update core logic %s", ref[:8]),
			Ref:    ref,
			Author: "direct_architect",
		}
	case "issue":
		return GitHubEvent{
			Type:   "repo.issue",
			Repo:   "lithix-src/catalyst",
			Title:  "bug: race condition in event bus",
			Ref:    ref,
			Author: "qa-bot",
		}
	default:
		return GitHubEvent{
			Type:   "repo.pr",
			Repo:   "lithix-src/catalyst",
			Title:  "chore: bump dependencies",
			Ref:    ref,
			Author: "dependabot",
		}
	}
}

func publishEvent(client mqtt.Client, evt GitHubEvent) {
	// Construct CloudEvent
	payload := map[string]interface{}{
		"id":          fmt.Sprintf("evt-%d", time.Now().UnixNano()),
		"source":      "repo-watcher",
		"specversion": "1.0",
		"type":        evt.Type,
		"time":        time.Now().UTC(),
		"data": map[string]string{
			"repo":   evt.Repo,
			"title":  evt.Title,
			"author": evt.Author,
			"ref":    evt.Ref,
		},
	}

	bytes, _ := json.Marshal(payload)
	topic := fmt.Sprintf("repo/%s/%s", "catalyst", "event") // Standardized topic for now

	token := client.Publish(topic, 0, false, bytes)
	token.Wait()
	log.Printf("[RepoWatcher] >>> NEW EVENT: %s | %s", evt.Type, evt.Title)
}
