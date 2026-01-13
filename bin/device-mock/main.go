package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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
	fmt.Println("Connected to Mosquitto. Starting Multi-Source Simulation...")

	// 1. Hardware Sensor Loop (Fast)
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for range ticker.C {
			temp := 45.0 + rand.Float64()*10.0
			emit(client, "sensor/cpu/temp", "sensor.cpu.temp", map[string]interface{}{"value": temp, "unit": "C"})
		}
	}()

	// 2. Agent Log Loop (Bursty)
	go func() {
		agents := []string{"TrendScout", "GapAnalyst", "CodeRunner"}
		actions := []string{"Scanning", "Analyzing", "Sleeping", "Fetching", "Compiling"}
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

	// 3. Repo Event Loop (Rare)
	go func() {
		for {
			time.Sleep(time.Duration(rand.Intn(10000)+5000) * time.Millisecond)
			issueID := rand.Intn(1000) + 1
			emit(client, "repo/lithix/issue", "repo.issue.new", map[string]interface{}{
				"repo":  "lithix/core",
				"id":    issueID,
				"title": fmt.Sprintf("Unexpected panic in worker %d", issueID),
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
	fmt.Printf("> Sent %s: %s\n", topic, eventType)
}
