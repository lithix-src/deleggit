package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/datacraft/catalyst/core/internal/adapter/llm"
	"github.com/datacraft/catalyst/core/internal/adapter/mqtt"
	"github.com/datacraft/catalyst/core/internal/adapter/store"
	"github.com/datacraft/catalyst/core/internal/domain"
	"github.com/datacraft/catalyst/core/internal/service"
	"github.com/point-unknown/catalyst/pkg/env"
	"github.com/point-unknown/catalyst/pkg/logger"
)

var (
	BrokerURL = env.Get("BROKER_URL", "tcp://localhost:1883")
	ClientID  = env.Get("CLIENT_ID", "catalyst-core-service")
	Log       *slog.Logger
)

func main() {
	Log = logger.New("core-service")
	Log.Info("Starting Catalyst Core Service (Phase 2)...")

	// ==========================================
	// 1. DOMAIN & SERVICE LAYER (The "Brain")
	// ==========================================

	// A. Registry
	registry := service.NewAgentRegistry()

	// 1.5 Data Store (PostgreSQL)
	// User NodePort 30000 -> 5432
	connStr := env.Get("DATABASE_URL", "postgres://catalyst:devpassword@localhost:5432/catalyst_core")
	pgStore, err := store.NewPostgresStore(connStr)
	if err != nil {
		Log.Warn("⚠️ [STORE] Failed to connect to Postgres (is K8s up?)", "error", err)
	} else {
		defer pgStore.Close()
		Log.Info("✅ [STORE] Connected to Postgres.")
		if err := pgStore.InitSchema(context.Background()); err != nil {
			Log.Error("Failed to init schema", "error", err)
			os.Exit(1)
		}
	}

	// 1.6 LLM Provider (The "Brain")
	llmCfg := llm.Config{
		Endpoint:    env.Get("LLM_ENDPOINT", "http://localhost:11434/v1"),
		Model:       env.Get("LLM_MODEL", "qwen2.5-coder:7b-instruct"),
		PrivateMode: env.Get("PRIVATE_MODE", "true") != "false",
	}

	llmProvider, err := llm.NewOpenAIAdapter(llmCfg)
	if err != nil {
		Log.Warn("⚠️ [LLM] Failed to init LLM Adapter. Thinking Disabled.", "error", err)
	} else {
		Log.Info("✅ [LLM] Connected", "endpoint", llmCfg.Endpoint)
	}

	// B. Register Standard Agents (Dynamic Loader)
	configFile := env.Get("CATALYST_CONFIG_PATH", "../../config/agents.yaml")

	loadedAgents, loadedMissions, err := service.LoadAgents(configFile, llmProvider)
	if err != nil {
		Log.Warn("⚠️ [LOADER] Failed to load agents.yaml. Running without agents.", "error", err)
	} else {
		for _, a := range loadedAgents {
			registry.Register(a)
			Log.Info("Registered Agent", "id", a.ID(), "type", a.Type())
		}
	}

	// C. Mission Manager (Orchestrator)
	missionMgr := service.NewMissionManager(registry)

	// D. Load Missions (Configuration)
	for _, m := range loadedMissions {
		missionMgr.LoadMission(m)
	}

	// ==========================================
	// 2. INFRASTRUCTURE LAYER (The "Plumbing")
	// ==========================================

	// A. EventBus (MQTT)
	mqttClient, err := mqtt.NewAdapter(BrokerURL, ClientID)
	if err != nil {
		Log.Error("Failed to connect to MQTT", "error", err)
		os.Exit(1)
	}
	defer mqttClient.Close()
	Log.Info("Connected to MQTT Broker")

	// B. Subscribe Routes
	// Route all sensor data to the Mission Manager
	err = mqttClient.Subscribe("sensor/#", func(event domain.CloudEvent) {
		// 1. Persist (Best Effort)
		if pgStore != nil {
			go func() {
				if err := pgStore.SaveEvent(context.Background(), event); err != nil {
					Log.Warn("Failed to save event", "error", err)
				}
			}()
		}
		// 2. Process
		missionMgr.ProcessEvent(event)
	})
	if err != nil {
		Log.Error("Failed to subscribe to sensors", "error", err)
		os.Exit(1)
	}

	Log.Info("Core Service Operational. Swarm is Active.")

	// Block until Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	Log.Info("Shutting down Core Service...")
}
