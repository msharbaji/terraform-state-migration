package backend

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/msharbaji/terraform-state-migration/terraform-hybrid/internal/config"
)

// Writer defines an interface for writing backend configuration
type Writer interface {
	WriteBackend(terraformConfig *config.TerraformHybridConfig, workspaceDir, callerName string) error
}

// TerraformBackendWriter implements Writer for different backends
type TerraformBackendWriter struct{}

// WriteBackend writes the backend configuration to the specified file
func (tbw *TerraformBackendWriter) WriteBackend(terraformConfig *config.TerraformHybridConfig, workspaceDir, callerName string) error {
	// Get the relative path under "deploy/provider"
	relativePath, err := tbw.getRelativePathUnderProvider(workspaceDir)
	if err != nil {
		return fmt.Errorf("error determining relative path: %v", err)
	}

	backendFile := filepath.Join(workspaceDir, "backend.tf")

	content, err := tbw.generateBackendContent(terraformConfig.Global.Backend, terraformConfig.Global.BackendType, relativePath)
	if err != nil {
		return fmt.Errorf("error generating backend content: %v", err)
	}

	if err := os.WriteFile(backendFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("error writing backend file %s: %v", backendFile, err)
	}

	fmt.Printf("Successfully wrote backend configuration to %s\n", backendFile)
	return nil
}

// getRelativePathUnderProvider calculates the relative path under "deploy/provider"
func (tbw *TerraformBackendWriter) getRelativePathUnderProvider(workspaceDir string) (string, error) {
	// Look for the "deploy/provider" folder in the absolute workspace path
	providerRoot := filepath.Join("deploy", "provider")
	absWorkspaceDir, err := filepath.Abs(workspaceDir)
	if err != nil {
		return "", fmt.Errorf("could not determine absolute path: %v", err)
	}

	// Find the index of "deploy/provider" in the absolute path
	index := strings.Index(absWorkspaceDir, providerRoot)
	if index == -1 {
		return "", fmt.Errorf("workspace directory does not seem to be under 'deploy/provider'")
	}

	// Return the relative path under "deploy/provider"
	return absWorkspaceDir[index+len(providerRoot)+1:], nil
}

// generateBackendContent generates the backend configuration content based on the backend type
func (tbw *TerraformBackendWriter) generateBackendContent(backend interface{}, backendType config.BackendType, relativePath string) (string, error) {
	switch backendType {
	case config.LocalBackendType:
		return tbw.generateLocalBackendContent(backend.(*config.LocalBackendConfig), relativePath)
	case config.BackendTypeCloudStorage:
		return tbw.generateCloudStorageBackendContent(backend.(*config.CloudStorageBackendConfig), relativePath)
	case config.BackendTypePostgres:
		return tbw.generatePostgresBackendContent(backend.(*config.PostgresBackendConfig))
	default:
		return "", fmt.Errorf("unsupported backend type: %s", backendType)
	}
}

// Local Backend
func (tbw *TerraformBackendWriter) generateLocalBackendContent(backend *config.LocalBackendConfig, relativePath string) (string, error) {
	// Use the relative path instead of subfolder name
	backendPath := fmt.Sprintf("%s/%s/terraform.tfstate", backend.Path, relativePath)
	return fmt.Sprintf(`terraform {
  backend "local" {
    path = "%s"
  }
}`, backendPath), nil
}

// Cloud Storage Backend
func (tbw *TerraformBackendWriter) generateCloudStorageBackendContent(backend *config.CloudStorageBackendConfig, relativePath string) (string, error) {
	bucketKey := fmt.Sprintf("%s/terraform.tfstate", relativePath)
	content := fmt.Sprintf(`terraform {
  backend "%s" {
    encrypt = "true"
    region  = "%s"
    bucket  = "%s"
    key     = "%s"
  }
}`, backend.Type, backend.Region, backend.BucketName, bucketKey)

	// Handle optional fields
	if backend.Endpoint != "" {
		content += fmt.Sprintf(`
    endpoint                    = "%s"
    skip_region_validation      = true
    skip_credentials_validation = true
    skip_metadata_api_check     = true`, backend.Endpoint)
	}

	if backend.RoleArn != "" {
		content += fmt.Sprintf(`
    role_arn  = "%s"`, backend.RoleArn)
	}

	content += "\n}"
	return content, nil
}

// Postgres Backend
func (tbw *TerraformBackendWriter) generatePostgresBackendContent(backend *config.PostgresBackendConfig) (string, error) {
	return fmt.Sprintf(`terraform {
  backend "pg" {
    conn_str     = "%s"
    schema_name  = "%s"
  }
}`, backend.ConnectionString, backend.SchemaName), nil
}
