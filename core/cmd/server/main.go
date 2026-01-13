package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/datacraft/deleggit/core/internal/adapter/agent"
	"github.com/datacraft/deleggit/core/internal/adapter/mqtt"
	"github.com/datacraft/deleggit/core/internal/domain"
	"github.com/datacraft/deleggit/core/internal/service"
)

const (
	BrokerURL = "tcp://localhost:1883"
	ClientID  = "deleggit-core-service"
)

func main() {
	log.Println("Starting Deleggit Core Service (Phase 2)...")

	// ==========================================
	// 1. DOMAIN & SERVICE LAYER (The "Brain")
	// ==========================================

	// A. Registry
	registry := service.NewAgentRegistry()

	// B. Register Standard Agents (Plugins)
	// In Phase 3 this will be dynamic. For now, we hardcode the "Swarm".
	registry.Register(agent.NewTrendScout("TrendScout", 80.0)) // Alert if avg > 80C
	registry.Register(agent.NewConsoleReporter("GapAnalyst"))

	// C. Mission Manager (Orchestrator)
	missionMgr := service.NewMissionManager(registry)

	// D. Load Default Missions (Configuration)
	// "When Sensor Data arrives, wake up TrendScout"
	missionMgr.LoadMission(domain.Mission{
		ID:           "mission-001",
		Name:         "Hardware Telemetry Watch",
		TriggerTopic: "sensor/#", // Matches any sensor event
		Agents:       []string{"TrendScout"},
	})

	// ==========================================
	// 2. INFRASTRUCTURE LAYER (The "Plumbing")
	// ==========================================

	// A. EventBus (MQTT)
	mqttClient, err := mqtt.NewAdapter(BrokerURL, ClientID)
	if err != nil {
		log.Fatalf("Failed to connect to MQTT: %v", err)
	}
	defer mqttClient.Close()
	log.Println("Connected to MQTT Broker")

	// B. Subscribe Routes
	// Route all sensor data to the Mission Manager
	err = mqttClient.Subscribe("sensor/#", func(event domain.CloudEvent) {
		missionMgr.ProcessEvent(event)
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to sensors: %v", err)
	}

	log.Println("Core Service Operational. Swarm is Active.")

	// Block until Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down Core Service...")
}
