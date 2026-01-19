package main

import (
	"context"
	"encoding/json"
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
	"github.com/point-unknown/catalyst/pkg/mcp"
	"github.com/point-unknown/catalyst/pkg/vector"
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
	// A. Registry
	agentRegistryService := service.NewAgentRegistry()

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

	// 1.7 MCP Registry (The "Hands")
	registry := mcp.NewLocalRegistry()
	registry.Register(mcp.Tool{
		Name:        "git_create_issue",
		Description: "Create a new issue in the repository",
		Parameters:  json.RawMessage(`{"type": "object", "properties": {"title": {"type": "string"}, "body": {"type": "string"}}, "required": ["title", "body"]}`),
	})
	registry.Register(mcp.Tool{
		Name:        "pipeline_list",
		Description: "List available CI/CD pipelines",
		Parameters:  json.RawMessage(`{"type": "object", "properties": {}, "required": []}`),
	})
	registry.Register(mcp.Tool{
		Name:        "pipeline_run",
		Description: "Trigger a CI/CD pipeline",
		Parameters:  json.RawMessage(`{"type": "object", "properties": {"workflow": {"type": "string"}}, "required": ["workflow"]}`),
	})
	Log.Info("✅ [MCP] Local Registry Initialized", "tools", 1)

	// B. Register Standard Agents (Dynamic Loader)
	configFile := env.Get("CATALYST_CONFIG_PATH", "config/agents.yaml")

	// Pass vector store to loader
	loadedAgents, loadedMissions, err := service.LoadAgents(configFile, llmProvider, registry, pgStore.Vector)
	if err != nil {
		Log.Warn("⚠️ [LOADER] Failed to load agents.yaml. Running without agents.", "error", err)
	} else {
		for _, a := range loadedAgents {
			// Register agent to the AgentRegistry Service
			agentRegistryService.Register(a)
			Log.Info("Registered Agent", "id", a.ID(), "type", a.Type())
		}
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

	// C. Mission Manager (Orchestrator)
	// Wrapper for MQTT Publish to match signature
	publisher := func(topic string, event domain.CloudEvent) {
		if err := mqttClient.Publish(topic, event); err != nil {
			Log.Warn("Failed to publish event from MissionManager", "topic", topic, "error", err)
		}
	}
	missionMgr := service.NewMissionManager(agentRegistryService, publisher)

	// D. Load Missions (Configuration)
	for _, m := range loadedMissions {
		missionMgr.LoadMission(m)
	}

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

	// Route Repo Events (Indexing)
	err = mqttClient.Subscribe("repo/#", func(event domain.CloudEvent) {
		if event.Type == "repo.content" {
			Log.Info("🧠 Indexing Code...", "size", len(event.Data))
			if pgStore != nil && pgStore.Vector != nil {
				go func() {
					ctx := context.Background()

					var payload map[string]string
					if err := json.Unmarshal(event.Data, &payload); err != nil {
						Log.Warn("Failed to unmarshal repo.content", "error", err)
						return
					}

					content := payload["content"]
					path := payload["path"]

					if content == "" {
						return
					}

					// 1. Embed
					vec, err := llmProvider.Embed(ctx, content)
					if err != nil {
						Log.Warn("Failed to embed code", "error", err)
						return // Retry?
					}

					// 2. Store
					doc := vector.Document{
						ID:        path, // Use path as ID to overwrite on update
						Content:   content,
						Embedding: vec,
						Metadata: map[string]interface{}{
							"hash": payload["hash"],
							"type": "code",
						},
					}

					if err := pgStore.Vector.Upsert(ctx, []vector.Document{doc}); err != nil {
						Log.Warn("Failed to upsert vector", "error", err)
					} else {
						Log.Info("✅ Indexed Code", "path", path)
					}
				}()
			}
		} else {
			// Normal Repo Events (PR, Push) -> Mission Manager
			missionMgr.ProcessEvent(event)
		}
	})

	Log.Info("Core Service Operational. Swarm is Active.")

	// Block until Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	Log.Info("Shutting down Core Service...")
}
