package workspace

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/datacraft/catalyst/core/internal/domain"
)

type LocalWorkspace struct {
	config domain.SecurityConfig
}

func NewLocalWorkspace(config domain.SecurityConfig) *LocalWorkspace {
	return &LocalWorkspace{config: config}
}

func (w *LocalWorkspace) ReadFile(path string) ([]byte, error) {
	fullPath, err := w.validatePath(path)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(fullPath)
}

func (w *LocalWorkspace) WriteFile(path string, data []byte) error {
	if w.config.ReadOnly {
		return fmt.Errorf("security policy violation: read-only mode")
	}

	fullPath, err := w.validatePath(path)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullPath, data, 0644)
}

func (w *LocalWorkspace) List(path string) ([]string, error) {
	fullPath, err := w.validatePath(path)
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
func (w *LocalWorkspace) validatePath(path string) (string, error) {
	// 1. Resolve Absolute Path
	// Handle relative paths from Sandbox
	cleanPath := filepath.Clean(path)
	// If it doesn't start with the sandbox, join it
	if !strings.HasPrefix(cleanPath, w.config.SandboxPath) {
		cleanPath = filepath.Join(w.config.SandboxPath, cleanPath)
	}

	// 2. Sandbox Escape Check
	if !strings.HasPrefix(cleanPath, w.config.SandboxPath) {
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
