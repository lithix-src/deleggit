package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	brokerURL  = "tcp://localhost:1883"
	clientID   = "catalyst-repo-watcher"
	repoPath   = "../../"              // Assuming running from bin/repo-watcher or root via make
	remoteName = "lithix-src/deleggit" // Default, can be parsed from git remote
)

// GitHubEvent structure
type GitHubEvent struct {
	Type   string
	Repo   string
	Title  string
	Ref    string
	Author string
	URL    string
}

func main() {
	log.Println("[RepoWatcher] Starting Real Git Poller...")

	// 1. MQTT Setup
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID(clientID)
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("[RepoWatcher] Failed to connect: %v", token.Error())
	}
	log.Println("[RepoWatcher] Connected to Event Bus")

	// 2. Initial State
	lastHash := getGitHeadHash()
	log.Printf("[RepoWatcher] Monitoring from Hash: %s", lastHash)

	// 3. Polling Loop
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				currentHash := getGitHeadHash()
				if currentHash != "" && currentHash != lastHash {
					log.Printf("[RepoWatcher] Change Detected! %s -> %s", lastHash, currentHash)

					// Fetch Commit Details
					details := getCommitDetails(currentHash)
					evt := GitHubEvent{
						Type:   "repo.push",
						Repo:   remoteName,
						Title:  details.Message,
						Ref:    currentHash,
						Author: details.Author,
						URL:    fmt.Sprintf("https://github.com/%s/commit/%s", remoteName, currentHash),
					}

					publishEvent(client, evt)
					lastHash = currentHash
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

func getGitHeadHash() string {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Error getting git head: %v", err)
		return ""
	}
	return strings.TrimSpace(string(out))
}

type CommitDetails struct {
	Author  string
	Message string
}

func getCommitDetails(hash string) CommitDetails {
	// git log -1 --format="%an|%s" hash
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

func publishEvent(client mqtt.Client, evt GitHubEvent) {
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
			"url":    evt.URL,
		},
	}

	bytes, _ := json.Marshal(payload)
	topic := fmt.Sprintf("repo/%s/%s", "catalyst", "event")

	token := client.Publish(topic, 0, false, bytes)
	token.Wait()
	log.Printf("[RepoWatcher] >>> NEW REAL CHANGE: %s", evt.Title)
}
