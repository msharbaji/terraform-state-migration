package backend

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/msharbaji/terraform-state-migration/terraform-hybrid/internal/config"
	"github.com/msharbaji/terraform-state-migration/terraform-hybrid/internal/utils"
)

// TerraformBackendManager handles the core logic for generating backend configurations
type TerraformBackendManager struct {
	configLoader   config.Loader
	folderFinder   utils.FolderFinder
	backendFactory WriteFactory
}

// NewTerraformBackendManager creates a new TerraformBackendManager instance
func NewTerraformBackendManager(
	configLoader config.Loader,
	folderFinder utils.FolderFinder,
	backendFactory WriteFactory,
) *TerraformBackendManager {
	return &TerraformBackendManager{
		configLoader:   configLoader,
		folderFinder:   folderFinder,
		backendFactory: backendFactory,
	}
}

// GenerateBackends orchestrates the loading of config, finding folders, and generating backend.tf files
func (tbm *TerraformBackendManager) GenerateBackends(configPath, providerFolderPath string) error {
	// Load the configuration
	loadedConfig, err := tbm.configLoader.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	// Determine provider folder based on config path (e.g., gcp, aws, ali)
	providerFolder := getProviderFolder(providerFolderPath, configPath)

	// Find the component folders under the provider folder
	componentFolders, err := tbm.folderFinder.FindComponentProviderFolders(providerFolder)
	if err != nil {
		return fmt.Errorf("error finding component provider folders: %v", err)
	}

	if len(loadedConfig.Global.Accounts) == 0 {
		// Walk through all component folders and their subdirectories
		for _, componentFolder := range componentFolders {
			if err := tbm.walkAndProcessComponentFolders(loadedConfig, componentFolder); err != nil {
				return fmt.Errorf("error processing folder %s: %v", componentFolder, err)
			}
		}
	} else {
		// Process only folders that match the account names
		for _, componentFolder := range componentFolders {
			if tbm.isFolderForAccount(componentFolder, loadedConfig.Global.Accounts) {
				if err := tbm.walkAndProcessComponentFolders(loadedConfig, componentFolder); err != nil {
					return fmt.Errorf("error processing folder %s: %v", componentFolder, err)
				}
			}
		}
	}

	return nil
}

// walkAndProcessComponentFolders walks the component folder and generates backend.tf in subdirectories
func (tbm *TerraformBackendManager) walkAndProcessComponentFolders(loadedConfig *config.TerraformHybridConfig, componentFolder string) error {
	// Use WalkDir to traverse all directories under the component folder
	err := filepath.WalkDir(componentFolder, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip if it's not a directory
		if !d.IsDir() {
			return nil
		}

		// Skip the component folder itself
		if filepath.Base(componentFolder) == filepath.Base(path) {
			return nil
		}

		// Skip the .terraform directory
		if d.Name() == ".terraform" || d.Name() == "terraform.tfstate.d" {
			return filepath.SkipDir
		}

		// Process the subfolder to generate backend.tf
		fmt.Printf("Processing subfolder: %s\n", path)
		if err := tbm.processFolder(loadedConfig, path); err != nil {
			return fmt.Errorf("error processing folder %s: %v", path, err)
		}

		return nil
	})

	return err
}

// Helper to get the provider folder based on config file name
func getProviderFolder(basePath, configPath string) string {
	configFile := filepath.Base(configPath)
	provider := strings.TrimSuffix(configFile, filepath.Ext(configFile))
	return filepath.Join(basePath, provider)
}

// Helper to check if folder matches account from config
func (tbm *TerraformBackendManager) isFolderForAccount(folder string, accounts map[string]string) bool {
	for account := range accounts {
		if strings.Contains(folder, account) {
			return true
		}
	}
	return false
}

// processFolder handles backend.tf generation for a specific folder
func (tbm *TerraformBackendManager) processFolder(loadedConfig *config.TerraformHybridConfig, folder string) error {
	writer, err := tbm.backendFactory.CreateBackendWriter(loadedConfig.Global.BackendType)
	if err != nil {
		return fmt.Errorf("error creating backend writer: %v", err)
	}

	// Write the backend configuration for the folder
	if err := writer.WriteBackend(loadedConfig, folder, "main"); err != nil {
		return fmt.Errorf("error writing backend for folder %s: %v", folder, err)
	}

	fmt.Printf("Successfully wrote backend for folder: %s\n", folder)
	return nil
}
