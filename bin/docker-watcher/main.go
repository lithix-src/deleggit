package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config
const (
	BrokerURL   = "tcp://localhost:1883"
	ClientID    = "catalyst-docker-watcher"
	TopicState  = "infra/docker/state"
	MetricsPort = ":8084"
)

var (
	containerCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "docker_containers_total",
		Help: "Total number of containers",
	})
	containerRunning = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "docker_containers_running",
		Help: "Number of running containers",
	})
)

func init() {
	prometheus.MustRegister(containerCount)
	prometheus.MustRegister(containerRunning)
}

// ContainerState represents the simplified view for the Dashboard
type ContainerState struct {
	ID     string  `json:"id"`
	Names  string  `json:"names"`
	Image  string  `json:"image"`
	State  string  `json:"state"`  // "running", "exited"
	Status string  `json:"status"` // "Up 2 hours"
	CPU    float64 `json:"cpu"`    // Placeholder for now
	Memory float64 `json:"memory"` // Placeholder for now
}

func main() {
	log.Println("[DockerWatcher] Starting...")

	// 1. Connect to Docker
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.44"))
	if err != nil {
		log.Fatalf("[DockerWatcher] Failed to create Docker client: %v", err)
	}
	defer cli.Close()

	// 2. Connect to MQTT
	opts := mqtt.NewClientOptions().AddBroker(BrokerURL).SetClientID(ClientID)
	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("[DockerWatcher] Failed to connect to MQTT: %v", token.Error())
	}
	log.Println("[DockerWatcher] Connected to MQTT Broker")

	// 3. Start Metrics Server
	go func() {
		log.Printf("[DockerWatcher] Serving metrics on %s", MetricsPort)
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(MetricsPort, nil))
	}()

	// 4. Polling Loop (Every 5s)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		pollContainers(cli, mqttClient)
	}
}

func pollContainers(cli *client.Client, m mqtt.Client) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		log.Printf("[DockerWatcher] Error listing containers: %v", err)
		return
	}

	runningCount := 0
	var states []ContainerState
	for _, c := range containers {
		name := "unknown"
		if len(c.Names) > 0 {
			name = c.Names[0]
		}

		if c.State == "running" {
			runningCount++
		}

		state := ContainerState{
			ID:     c.ID[:12],
			Names:  name,
			Image:  c.Image,
			State:  c.State,
			Status: c.Status,
			CPU:    0.0, // Future: Calculate from Stats API
			Memory: 0.0, // Future: Calculate from Stats API
		}
		states = append(states, state)
	}

	// Update Metrics
	containerCount.Set(float64(len(containers)))
	containerRunning.Set(float64(runningCount))

	// Publish List (CloudEvent Format)
	payload := map[string]interface{}{
		"id":          fmt.Sprintf("evt-%d", time.Now().UnixNano()),
		"source":      "docker-watcher",
		"specversion": "1.0",
		"type":        "infra.docker.state",
		"time":        time.Now().UTC(),
		"data":        states,
	}

	bytes, _ := json.Marshal(payload)
	m.Publish(TopicState, 0, false, bytes)

	// Publish CSS Compliant Sensor Data (Dynamic Grid)
	// 1. Running Count
	sensorPayload := map[string]interface{}{
		"value": float64(runningCount),
		"unit":  "cnt",
		"label": "Running Containers",
		"meta": map[string]string{
			"total": fmt.Sprintf("%d", len(containers)),
		},
	}
	sensorBytes, _ := json.Marshal(map[string]interface{}{
		"id":          fmt.Sprintf("evt-%d", time.Now().UnixNano()),
		"source":      "docker-watcher",
		"specversion": "1.0",
		"type":        "sensor.docker.running",
		"time":        time.Now().UTC(),
		"data":        sensorPayload,
	})
	m.Publish("sensor/docker/running", 0, false, sensorBytes)

	// Also log to swarm activity occasionally
	if len(states) > 0 {
		logPayload := fmt.Sprintf(`{"agent": "Infrastructure", "message": "Monitoring %d containers (%d running)"}`, len(states), runningCount)
		m.Publish("agent/infra/log", 0, false, logPayload)
	}
}
