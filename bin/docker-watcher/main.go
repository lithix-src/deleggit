package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/point-unknown/catalyst/pkg/cloudevent"
	"github.com/point-unknown/catalyst/pkg/env"
	"github.com/point-unknown/catalyst/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config via SDK
var (
	BrokerURL   = env.Get("BROKER_URL", "tcp://localhost:1883")
	ClientID    = env.Get("CLIENT_ID", "catalyst-docker-watcher")
	MetricsPort = env.Get("METRICS_PORT", ":8084")
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
	log := logger.New("docker-watcher")
	log.Info("Starting Docker Watcher...")

	// 1. Connect to Docker
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.44"))
	if err != nil {
		log.Error("Failed to create Docker client", "error", err)
		return
	}
	defer cli.Close()

	// 2. Connect to MQTT
	opts := mqtt.NewClientOptions().AddBroker(BrokerURL).SetClientID(ClientID)
	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Error("Failed to connect to MQTT", "error", token.Error())
		return
	}
	log.Info("Connected to MQTT Broker", "url", BrokerURL)

	// 3. Start Metrics Server
	go func() {
		log.Info("Serving metrics", "port", MetricsPort)
		http.Handle("/metrics", promhttp.Handler())
		// http.ListenAndServe is blocking, so we panic if it fails to ensure visibility
		panic(http.ListenAndServe(MetricsPort, nil))
	}()

	// 4. Polling Loop (Every 5s)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		pollContainers(cli, mqttClient, log)
	}
}

func pollContainers(cli *client.Client, m mqtt.Client, log *slog.Logger) {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		log.Error("Error listing containers", "error", err)
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
			CPU:    0.0,
			Memory: 0.0,
		}
		states = append(states, state)
	}

	// Update Metrics
	containerCount.Set(float64(len(containers)))
	containerRunning.Set(float64(runningCount))

	// 1. Publish State List (Data Dump)
	publishEvent(m, log, "infra/docker/state", "infra.docker.state", states)

	// 2. Publish CSS Compliant Sensor Data (Running Count)
	sensorPayload := map[string]interface{}{
		"value": float64(runningCount),
		"unit":  "cnt",
		"label": "Running Containers",
		"meta": map[string]string{
			"total": fmt.Sprintf("%d", len(containers)),
		},
	}
	publishEvent(m, log, "sensor/docker/running", "sensor.docker.running", sensorPayload)

	log.Info("Polled Docker", "total", len(containers), "running", runningCount)
}

func publishEvent(client mqtt.Client, log *slog.Logger, topic, eventType string, data interface{}) {
	evt, err := cloudevent.New("docker-watcher", eventType, data)
	if err != nil {
		log.Error("Failed to create event", "type", eventType, "error", err)
		return
	}

	bytes, err := json.Marshal(evt)
	if err != nil {
		log.Error("Failed to marshal event", "type", eventType, "error", err)
		return
	}

	token := client.Publish(topic, 0, false, bytes)
	token.Wait()
}
