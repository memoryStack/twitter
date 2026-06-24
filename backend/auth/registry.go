package auth

import (
	"fmt"
	"sync"
)

// Registry holds all registered authentication strategies.
// Use Register during init; handlers resolve the active strategy by name.
var Registry = &strategyRegistry{
	byName: make(map[string]Strategy),
}

type strategyRegistry struct {
	mu          sync.RWMutex
	byName      map[string]Strategy
	defaultName string
}

// Register adds a strategy. Panics if the name is already registered.
func (r *strategyRegistry) Register(s Strategy) {
	if s == nil {
		panic("auth: cannot register nil strategy")
	}
	name := s.Name()
	if name == "" {
		panic("auth: strategy name is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.byName[name]; exists {
		panic(fmt.Sprintf("auth: strategy %q already registered", name))
	}
	r.byName[name] = s
	if r.defaultName == "" {
		r.defaultName = name
	}
}

// SetDefault marks which strategy is used when ?mode= is omitted.
func (r *strategyRegistry) SetDefault(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byName[name]; !ok {
		return fmt.Errorf("auth: unknown strategy %q", name)
	}
	r.defaultName = name
	return nil
}

// Get returns a strategy by name. Empty name resolves to the default strategy.
func (r *strategyRegistry) Get(name string) (Strategy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if name == "" {
		name = r.defaultName
	}
	if name == "" {
		return nil, fmt.Errorf("auth: no strategies registered")
	}
	s, ok := r.byName[name]
	if !ok {
		return nil, fmt.Errorf("auth: unknown strategy %q", name)
	}
	return s, nil
}

// Default returns the default strategy.
func (r *strategyRegistry) Default() (Strategy, error) {
	return r.Get("")
}

// All returns every registered strategy (useful for token validation across modes).
func (r *strategyRegistry) All() []Strategy {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Strategy, 0, len(r.byName))
	for _, s := range r.byName {
		out = append(out, s)
	}
	return out
}
