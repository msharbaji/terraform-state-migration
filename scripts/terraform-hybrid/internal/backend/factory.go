package backend

import (
	"fmt"

	"github.com/msharbaji/terraform-state-migration/terraform-hybrid/internal/config"
)

// WriteFactory is responsible for creating backend writers based on the backend type
type WriteFactory struct{}

func NewBackendFactory() *WriteFactory {
	return &WriteFactory{}
}

// CreateBackendWriter creates a backend writer based on the given backend type
func (f *WriteFactory) CreateBackendWriter(backendType config.BackendType) (Writer, error) {
	switch backendType {
	case config.LocalBackendType:
		return &TerraformBackendWriter{}, nil
	case config.BackendTypeCloudStorage:
		return &TerraformBackendWriter{}, nil
	case config.BackendTypePostgres:
		return &TerraformBackendWriter{}, nil
	default:
		return nil, fmt.Errorf("unsupported backend type: %s", backendType)
	}
}
