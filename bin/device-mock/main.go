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
			// Get total CPU usage
			c, err := cpu.Percent(0, false)
			if err == nil && len(c) > 0 {
				// Emit Usage
				emit(client, "sensor/cpu/usage", "sensor.cpu.usage", map[string]interface{}{
					"value": c[0],
					"unit":  "%",
					"label": "CPU Load",
					"meta":  map[string]string{"status": "nominal"},
				})
			}
		}
	}()

	// 1.5 Real Memory Sensor Loop
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for range ticker.C {
			v, err := mem.VirtualMemory()
			if err == nil {
				// Log to console to debug 0.0 issue
				// fmt.Printf("Mem: Total=%v, Used=%v, Percent=%v\n", v.Total, v.Used, v.UsedPercent)
				emit(client, "sensor/memory/usage", "sensor.memory.usage", map[string]interface{}{
					"value": v.UsedPercent,
					"unit":  "%",
					"label": "Memory Usage",
				})
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

	// 3. Repo Event Loop - REMOVED (Handled by dedicated repo-watcher service)

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
