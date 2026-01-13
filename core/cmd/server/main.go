package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/datacraft/catalyst/core/internal/adapter/agent"
	"github.com/datacraft/catalyst/core/internal/adapter/mqtt"
	"github.com/datacraft/catalyst/core/internal/adapter/store"
	"github.com/datacraft/catalyst/core/internal/domain"
	"github.com/datacraft/catalyst/core/internal/service"
)

const (
	BrokerURL = "tcp://localhost:1883"
	ClientID  = "catalyst-core-service"
)

func main() {
	log.Println("Starting Catalyst Core Service (Phase 2)...")

	// ==========================================
	// 1. DOMAIN & SERVICE LAYER (The "Brain")
	// ==========================================

	// A. Registry
	registry := service.NewAgentRegistry()

	// 1.5 Data Store (PostgreSQL)
	// User NodePort 30000 -> 5432
	connStr := "postgres://catalyst:devpassword@localhost:5432/catalyst_core"
	pgStore, err := store.NewPostgresStore(connStr)
	if err != nil {
		log.Printf("⚠️ [STORE] Failed to connect to Postgres (is K8s up?): %v", err)
	} else {
		defer pgStore.Close()
		log.Println("✅ [STORE] Connected to Postgres.")
		if err := pgStore.InitSchema(context.Background()); err != nil {
			log.Fatalf("Failed to init schema: %v", err)
		}
	}

	// B. Register Standard Agents (Plugins)
	// In Phase 3 this will be dynamic. For now, we hardcode the "Swarm".
	registry.Register(agent.NewTrendScout("TrendScout", 80.0)) // Alert if avg > 80C
	registry.Register(agent.NewConsoleReporter("SwarmLog"))    // General Logger

	// C. Mission Manager (Orchestrator)
	missionMgr := service.NewMissionManager(registry)

	// D. Load Default Missions (Configuration)
	// 1. Hardware Watch
	missionMgr.LoadMission(domain.Mission{
		ID:           "mission-001",
		Name:         "Hardware Telemetry Watch",
		TriggerTopic: "sensor/#", // Matches any sensor event
		Agents:       []string{"TrendScout"},
	})

	// 2. Swarm Activity Watch
	missionMgr.LoadMission(domain.Mission{
		ID:           "mission-002",
		Name:         "Swarm Activity Stream",
		TriggerTopic: "agent/+/log", // Matches agent logs
		Agents:       []string{"SwarmLog"},
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
		// 1. Persist (Best Effort)
		if pgStore != nil {
			go func() {
				if err := pgStore.SaveEvent(context.Background(), event); err != nil {
					log.Printf("Failed to save event: %v", err)
				}
			}()
		}
		// 2. Process
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
