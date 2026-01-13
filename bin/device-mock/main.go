package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// SensorData adheres to the CloudEvents-like structure
type CloudEvent struct {
	ID        string      `json:"id"`
	Source    string      `json:"source"`
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"time"`
}

func main() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("device-mock-multi")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Connected to Mosquitto. Starting Real Hardware Telemetry...")

	// 1. Real CPU Sensor Loop
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for range ticker.C {
			// Get total CPU usage (false = all cores aggregated)
			c, err := cpu.Percent(0, false)
			if err == nil && len(c) > 0 {
				// Note: using 'value' key to match UI expectation
				// UI expects Key "value" for the chart
				emit(client, "sensor/cpu/temp", "sensor.cpu.usage", map[string]interface{}{"value": c[0], "unit": "%"})
			}
		}
	}()

	// 1.5 Real Memory Sensor Loop
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for range ticker.C {
			v, err := mem.VirtualMemory()
			if err == nil {
				emit(client, "sensor/memory/usage", "sensor.memory.usage", map[string]interface{}{"value": v.UsedPercent, "unit": "%"})
			}
		}
	}()

	// 2. Agent Log Loop (Still Simulated for now)
	go func() {
		agents := []string{"Interface", "Orchestrator", "Infrastructure", "Compliance", "Simulation"}
		actions := []string{"Optimizing", "Compiling", "Provisioning", "Verifying", "Simulating"}
		for {
			time.Sleep(time.Duration(rand.Intn(3000)+500) * time.Millisecond)
			agent := agents[rand.Intn(len(agents))]
			action := actions[rand.Intn(len(actions))]
			msg := fmt.Sprintf("[%s] %s target...", agent, action)

			emit(client, fmt.Sprintf("agent/%s/log", agent), "agent.log", map[string]interface{}{
				"agent":   agent,
				"level":   "INFO",
				"message": msg,
			})
		}
	}()

	// 3. Repo Event Loop (Simulated Bridge to Real URLs)
	go func() {
		repos := []string{"catalyst/ui", "catalyst/core", "catalyst/infra"}
		for {
			time.Sleep(time.Duration(rand.Intn(10000)+5000) * time.Millisecond)
			repo := repos[rand.Intn(len(repos))]
			issueID := rand.Intn(1000) + 1
			emit(client, "repo/issue/new", "repo.issue.new", map[string]interface{}{
				"repo":  repo,
				"id":    issueID,
				"title": fmt.Sprintf("Anomaly detected in %s module", repo),
			})
		}
	}()

	select {} // Block forever
}

func emit(client mqtt.Client, topic string, eventType string, data interface{}) {
	payload := CloudEvent{
		ID:        fmt.Sprintf("evt-%d", time.Now().UnixNano()),
		Source:    "device-mock",
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now(),
	}
	bytes, _ := json.Marshal(payload)
	token := client.Publish(topic, 0, false, bytes)
	token.Wait()
	// fmt.Printf("> Sent %s: %s\n", topic, eventType) // Reduce spam
}
