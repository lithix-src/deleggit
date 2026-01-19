package mcp

import "sync"

// LocalRegistry is a thread-safe in-memory tool registry.
type LocalRegistry struct {
	tools map[string]Tool
	mu    sync.RWMutex
}

func NewLocalRegistry() *LocalRegistry {
	return &LocalRegistry{
		tools: make(map[string]Tool),
	}
}

func (r *LocalRegistry) Register(t Tool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools[t.Name] = t
}

func (r *LocalRegistry) ListTools() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]Tool, 0, len(r.tools))
	for _, t := range r.tools {
		list = append(list, t)
	}
	return list
}

func (r *LocalRegistry) GetTool(name string) (*Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tools[name]
	return &t, ok
}
