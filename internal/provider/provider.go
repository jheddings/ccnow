package provider

import "github.com/jheddings/ccglow/internal/types"

// RegisterBuiltin adds all built-in data providers to the registry.
func RegisterBuiltin(registry *Registry) {
	registry.Register(&pwdProvider{})
	registry.Register(&gitProvider{})
	registry.Register(&contextProvider{})
	registry.Register(&modelProvider{})
	registry.Register(&costProvider{})
	registry.Register(&sessionProvider{})
	registry.Register(&speedProvider{})
	registry.Register(&claudeProvider{})
	registry.Register(&systemProvider{})
}

// Registry maps provider names to their implementations.
type Registry struct {
	providers map[string]types.DataProvider
}

// NewRegistry creates an empty provider registry.
func NewRegistry() *Registry {
	return &Registry{providers: make(map[string]types.DataProvider)}
}

// Register adds a data provider implementation.
func (r *Registry) Register(p types.DataProvider) {
	r.providers[p.Name()] = p
}

// All returns the underlying provider map for use by the render pipeline.
func (r *Registry) All() map[string]types.DataProvider {
	return r.providers
}
