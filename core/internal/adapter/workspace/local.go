package workspace

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/datacraft/catalyst/core/internal/domain"
)

type LocalWorkspace struct {
	config   domain.SecurityConfig
	resolver domain.ContextResolver
}

func NewLocalWorkspace(config domain.SecurityConfig, resolver domain.ContextResolver) *LocalWorkspace {
	return &LocalWorkspace{config: config, resolver: resolver}
}

func (w *LocalWorkspace) ReadFile(ctx context.Context, path string) ([]byte, error) {
	fullPath, err := w.validatePath(ctx, path)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(fullPath)
}

func (w *LocalWorkspace) WriteFile(ctx context.Context, path string, data []byte) error {
	if w.config.ReadOnly {
		return fmt.Errorf("security policy violation: read-only mode")
	}

	fullPath, err := w.validatePath(ctx, path)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullPath, data, 0644)
}

func (w *LocalWorkspace) List(ctx context.Context, path string) ([]string, error) {
	fullPath, err := w.validatePath(ctx, path)
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, f := range files {
		names = append(names, f.Name())
	}
	return names, nil
}

// validatePath ensures the path is inside the sandbox and allowed by patterns
func (w *LocalWorkspace) validatePath(ctx context.Context, path string) (string, error) {
	// 0. Resolve Sandbox (Dynamic)
	sandboxPath := w.config.SandboxPath

	// If Config has no static sandbox, OR if we have a resolver, use the resolver.
	// We prefer the resolver (Active Context) if available.
	if w.resolver != nil {
		activePath, err := w.resolver.GetActiveRepoPath(ctx)
		if err == nil && activePath != "" {
			sandboxPath = activePath
		}
		// If error (e.g. no active context), we might fallback to static or fail.
		// For now, let's log/fail? Or fallback to config.SandboxPath
	}

	if sandboxPath == "" {
		return "", fmt.Errorf("no active workspace context")
	}

	// 1. Resolve Absolute Path
	// Handle relative paths from Sandbox
	cleanPath := filepath.Clean(path)
	// If it doesn't start with the sandbox, join it
	if !strings.HasPrefix(cleanPath, sandboxPath) {
		cleanPath = filepath.Join(sandboxPath, cleanPath)
	}

	// 2. Sandbox Escape Check
	if !strings.HasPrefix(cleanPath, sandboxPath) {
		return "", fmt.Errorf("security violation: path traversal attempt (%s)", path)
	}

	// 3. Deny List Check
	base := filepath.Base(cleanPath)
	for _, pattern := range w.config.DeniedPatterns {
		matched, _ := filepath.Match(pattern, base)
		if matched {
			return "", fmt.Errorf("security violation: access denied to %s", pattern)
		}
	}

	// 4. Allow List Check (If defined)
	if len(w.config.AllowedPatterns) > 0 {
		allowed := false
		for _, pattern := range w.config.AllowedPatterns {
			matched, _ := filepath.Match(pattern, base)
			if matched {
				allowed = true
				break
			}
		}
		if !allowed {
			return "", fmt.Errorf("security violation: file type not allowed")
		}
	}

	return cleanPath, nil
}
