package web

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/datacraft/catalyst/core/internal/adapter/store"
	"github.com/datacraft/catalyst/core/internal/service"
	"github.com/point-unknown/catalyst/pkg/logger"
)

type Server struct {
	router    *http.ServeMux
	store     *store.PostgresStore
	workspace *service.WorkspaceManager
	log       *slog.Logger
}

func NewServer(store *store.PostgresStore, workspace *service.WorkspaceManager) *Server {
	s := &Server{
		router:    http.NewServeMux(),
		store:     store,
		workspace: workspace,
		log:       logger.New("web-adapter"),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	// CORS Middleware
	s.router.Handle("/api/agents", s.cors(http.HandlerFunc(s.handleAgents)))
	s.router.Handle("/api/repos", s.cors(http.HandlerFunc(s.handleRepos)))
	s.router.Handle("/api/context", s.cors(http.HandlerFunc(s.handleContext)))
}

func (s *Server) Run(addr string) error {
	s.log.Info("Starting HTTP API", "addr", addr)
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Handlers

func (s *Server) handleAgents(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := s.store.Pool().Query(ctx, "SELECT id, service, role, config FROM agents")
	if err != nil {
		s.log.Error("Failed to query agents", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var agents []map[string]interface{}
	for rows.Next() {
		var id, service, role string
		var config []byte
		if err := rows.Scan(&id, &service, &role, &config); err != nil {
			continue
		}
		agents = append(agents, map[string]interface{}{
			"id":      id,
			"service": service,
			"role":    role,
			"config":  json.RawMessage(config),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agents)
}

func (s *Server) handleRepos(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := s.store.Pool().Query(ctx, "SELECT id, org, name, default_branch FROM repos")
	if err != nil {
		s.log.Error("Failed to query repos", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var repos []map[string]interface{}
	for rows.Next() {
		var id, org, name, branch string
		if err := rows.Scan(&id, &org, &name, &branch); err != nil {
			continue
		}
		repos = append(repos, map[string]interface{}{
			"id":             id,
			"org":            org,
			"name":           name,
			"default_branch": branch,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
}

func (s *Server) handleContext(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if r.Method == "GET" {
		var activeRepoID, activeBranch string
		err := s.store.Pool().QueryRow(ctx, "SELECT active_repo_id, active_branch FROM context WHERE singleton_id = TRUE").Scan(&activeRepoID, &activeBranch)
		if err != nil {
			s.log.Error("Failed to get context", "error", err)
			http.Error(w, "Context not found (db init pending?)", http.StatusNotFound)
			return
		}

		// Hydrate Repo details
		var org, name string
		err = s.store.Pool().QueryRow(ctx, "SELECT org, name FROM repos WHERE id = $1", activeRepoID).Scan(&org, &name)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"active_repo_id": activeRepoID,
			"active_branch":  activeBranch,
			"org":            org,
			"name":           name,
		})
		return
	}

	if r.Method == "POST" {
		var payload struct {
			RepoID string `json:"repo_id"`
			Branch string `json:"branch"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid body", http.StatusBadRequest)
			return
		}

		// Use Logic Layer (Workspace Manager) to handle safe switching
		if err := s.workspace.SwitchContext(r.Context(), payload.RepoID, payload.Branch); err != nil {
			s.log.Error("Failed to switch context", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}
}
