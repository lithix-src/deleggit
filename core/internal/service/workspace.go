package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/datacraft/catalyst/core/internal/adapter/store"
	"github.com/point-unknown/catalyst/pkg/logger"
)

type WorkspaceManager struct {
	store *store.PostgresStore
	root  string // Absolute path to workspace root (d:\Datacraft)
	log   *slog.Logger
}

func NewWorkspaceManager(store *store.PostgresStore, rootPath string) *WorkspaceManager {
	// Check for env var override
	if envRoot := os.Getenv("WORKSPACE_ROOT"); envRoot != "" {
		rootPath = envRoot
	}

	abs, err := filepath.Abs(rootPath)
	if err != nil {
		abs = rootPath
	}
	return &WorkspaceManager{
		store: store,
		root:  abs,
		log:   logger.New("workspace-manager"),
	}
}

// SwitchContext verifies the target repo exists on disk, clones if missing (TODO), and updates DB.
func (w *WorkspaceManager) SwitchContext(ctx context.Context, repoID, branch string) error {
	w.log.Info("Switching Context", "repo", repoID, "branch", branch)

	// 1. Get Repo Details
	var name, org, defaultBranch string
	err := w.store.Pool().QueryRow(ctx, "SELECT name, org, default_branch FROM repos WHERE id = $1", repoID).Scan(&name, &org, &defaultBranch)
	if err != nil {
		return fmt.Errorf("repo not found: %w", err)
	}

	// 2. Resolve Path (/workspace/projects/<org>/<name>)
	// If root is /workspace, path is /workspace/projects/org/name
	targetPath := filepath.Join(w.root, "projects", org, name)

	// 3. Verify Existence
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		w.log.Info("Repo missing locally. Cloning...", "path", targetPath)
		// 3.1 Clone (Simple SSH clone)
		// git clone git@github.com:<org>/<name>.git <targetPath>
		repoURL := fmt.Sprintf("git@github.com:%s/%s.git", org, name)
		cmd := exec.CommandContext(ctx, "git", "clone", repoURL, targetPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to clone %s: %s", repoURL, string(out))
		}
		w.log.Info("✅ Clone Successful")
	}

	// 4. Update Context in DB
	_, err = w.store.Pool().Exec(ctx, "UPDATE context SET active_repo_id = $1, active_branch = $2 WHERE singleton_id = TRUE", repoID, branch)
	if err != nil {
		return fmt.Errorf("failed to update db context: %w", err)
	}

	w.log.Info("Context Switched Successfully", "path", targetPath)
	return nil
}

func (w *WorkspaceManager) GetWorkspaceRoot() string {
	return w.root
}

// GetActiveRepoPath returns the absolute path to the currently active repository
func (w *WorkspaceManager) GetActiveRepoPath(ctx context.Context) (string, error) {
	var activeRepoID string
	err := w.store.Pool().QueryRow(ctx, "SELECT active_repo_id FROM context WHERE singleton_id = TRUE").Scan(&activeRepoID)
	if err != nil {
		return "", fmt.Errorf("failed to get active context: %w", err)
	}

	var name, org string
	// We need 'org' now too
	err = w.store.Pool().QueryRow(ctx, "SELECT name, org FROM repos WHERE id = $1", activeRepoID).Scan(&name, &org)
	if err != nil {
		return "", fmt.Errorf("failed to get repo details: %w", err)
	}

	return filepath.Join(w.root, "projects", org, name), nil
}
