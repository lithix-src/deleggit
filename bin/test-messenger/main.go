package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Simplified CloudEvent for testing
type CloudEvent struct {
	ID     string          `json:"id"`
	Source string          `json:"source"`
	Type   string          `json:"type"`
	Time   time.Time       `json:"time"`
	Data   json.RawMessage `json:"data"`
}

func main() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("test-messenger-" + fmt.Sprint(time.Now().Unix()))

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect: %v", token.Error())
	}
	defer client.Disconnect(250)
	log.Println("✅ Connected to MQTT")

	// Subscribe to Agent Logs to see response
	client.Subscribe("agent/+/log", 0, func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("📩 RECEIVED [%s]: %s\n", msg.Topic(), string(msg.Payload()))
	})
	// Also subscribe to tool calls
	client.Subscribe("tool/call", 0, func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("🛠️ TOOL CALL [%s]: %s\n", msg.Topic(), string(msg.Payload()))
	})

	// Subscribe to EVERYTHING for debugging
	client.Subscribe("#", 0, func(client mqtt.Client, msg mqtt.Message) {
		if msg.Topic() == "agent/Liaison/wake" {
			return
		} // Ignore self
		log.Printf("📡 BUS [%s]: %s\n", msg.Topic(), string(msg.Payload()))
	})

	time.Sleep(1 * time.Second)

	// Send Message to Liaison
	payload := map[string]interface{}{
		"user_input": "Hello Agent Swarm, this is a bus test.",
	}
	data, _ := json.Marshal(payload)

	event := CloudEvent{
		ID:     fmt.Sprintf("evt-%d", time.Now().Unix()),
		Source: "test-messenger",
		Type:   "agent/Liaison/wake", // Standardized to slashes
		Time:   time.Now(),
		Data:   json.RawMessage(data),
	}
	// Check main.go: "agent.liaison.wake" (lowercase ID?)
	// In agents.yaml ID is "Liaison". In main.go check: fmt.Sprintf("agent.%s.wake", a.id).
	// If ID is "Liaison", it expects "agent.Liaison.wake".

	// Let's send to "agent.liaison.wake" AND "agent.Liaison.wake" to cover bases or check logging.
	// Actually, mission-chat TRIGGER dictates what TOPIC it listens to.

	topic := "agent/Liaison/wake"

	// We publish the CloudEvent JSON
	eventBytes, _ := json.Marshal(event)
	log.Printf("Tb Publishing to %s...\n", topic)
	token := client.Publish(topic, 0, false, eventBytes)
	token.Wait()
	if token.Error() != nil {
		log.Printf("❌ Publish Failed: %v\n", token.Error())
	} else {
		log.Println("✅ Publish Sent")
	}

	// Wait for response
	log.Println("Waiting for response...")
	time.Sleep(10 * time.Second)
}
