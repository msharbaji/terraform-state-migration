package commands

import (
	"fmt"

	"github.com/msharbaji/terraform-state-migration/terraform-hybrid/internal/backend"
	"github.com/msharbaji/terraform-state-migration/terraform-hybrid/internal/config"
	"github.com/msharbaji/terraform-state-migration/terraform-hybrid/internal/utils"
)

// GenerateBackendCmd defines the structure for the GenerateBackend command
type GenerateBackendCmd struct {
	Config         string `help:"Path to the YAML config file." required:"true" type:"path"`
	ProviderFolder string `help:"Path to the provider folder." required:"true" type:"path"`
}

// Run executes the logic for the GenerateBackend command
func (g *GenerateBackendCmd) Run() error {
	fmt.Printf("Using config file: %s\n", g.Config)
	fmt.Printf("Using provider folder: %s\n", g.ProviderFolder)

	// Initialize the config loader and utilities
	configLoader := config.NewConfigLoader()
	folderFinder := utils.NewFolderFinder()
	backendFactory := backend.NewBackendFactory()

	// Create the backend manager
	manager := backend.NewTerraformBackendManager(configLoader, folderFinder, *backendFactory)

	// Generate the backend configuration
	if err := manager.GenerateBackends(g.Config, g.ProviderFolder); err != nil {
		return fmt.Errorf("error generating backends: %w", err)
	}

	fmt.Println("Backend generation completed successfully.")
	return nil
}
