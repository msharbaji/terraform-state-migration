package config

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

var (
	DefaultConfigName = "aws.yaml"
)

type BackendType string

const (
	LocalBackendType        BackendType = "local"
	BackendTypeCloudStorage BackendType = "cloud_storage"
	BackendTypePostgres     BackendType = "postgres"
)

// String returns the string representation of the BackendType
func (bt BackendType) String() string {
	return string(bt)
}

// LocalBackendConfig represents the local backend configuration
type LocalBackendConfig struct {
	Path string `yaml:"path"`
}

// CloudStorageBackendConfig represents the configuration for cloud storage
type CloudStorageBackendConfig struct {
	Region     string `yaml:"region"`
	BucketName string `yaml:"bucket_name"`
	Type       string `yaml:"type"`
	RoleArn    string `yaml:"role_arn"`
	Endpoint   string `yaml:"endpoint"`
}

// PostgresBackendConfig represents the configuration for Postgres
type PostgresBackendConfig struct {
	ConnectionString string `yaml:"connection_string" validate:"required"`
	SchemaName       string `yaml:"schema_name" validate:"required"`
}

// GlobalConfig represents the global configuration
type GlobalConfig struct {
	BackendType BackendType       `yaml:"backend_type"`
	Backend     interface{}       `yaml:"-"`
	Accounts    map[string]string `yaml:"accounts"`
}

// TerraformHybridConfig represents the entire configuration
type TerraformHybridConfig struct {
	Global GlobalConfig `yaml:"global"`
}

// UnmarshalYAML unmarshals the YAML configuration into the GlobalConfig struct
func (gc *GlobalConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Temporary struct to hold common fields
	var temp struct {
		BackendType BackendType            `yaml:"backend_type"`
		Accounts    map[string]string      `yaml:"accounts"`
		Backend     map[string]interface{} `yaml:"backend"`
	}

	if err := unmarshal(&temp); err != nil {
		return err
	}

	gc.BackendType = temp.BackendType
	gc.Accounts = temp.Accounts

	switch temp.BackendType {
	case LocalBackendType:
		var local LocalBackendConfig
		data, err := yaml.Marshal(temp.Backend)
		if err != nil {
			return fmt.Errorf("error marshalling local backend: %v", err)
		}
		if err = yaml.Unmarshal(data, &local); err != nil {
			return fmt.Errorf("error unmarshalling local backend: %v", err)
		}
		gc.Backend = &local // Store as a pointer

	case BackendTypeCloudStorage:
		var cloudStorage CloudStorageBackendConfig
		data, err := yaml.Marshal(temp.Backend)
		if err != nil {
			return err
		}
		if err = yaml.Unmarshal(data, &cloudStorage); err != nil {
			return fmt.Errorf("error unmarshalling cloud storage backend: %v", err)
		}
		gc.Backend = &cloudStorage // Store as a pointer

	case BackendTypePostgres:
		var postgres PostgresBackendConfig
		data, err := yaml.Marshal(temp.Backend)
		if err != nil {
			return fmt.Errorf("error marshalling postgres backend: %v", err)
		}
		if err = yaml.Unmarshal(data, &postgres); err != nil {
			return fmt.Errorf("error unmarshalling postgres backend: %v", err)
		}
		gc.Backend = &postgres // Store as a pointer

	default:
		return fmt.Errorf("unknown backend_type: %s", temp.BackendType)
	}

	return nil
}

// CloudStorageBackend returns the CloudStorageBackendConfig from GlobalConfig
func (gc *GlobalConfig) CloudStorageBackend() (*CloudStorageBackendConfig, error) {
	if gc.BackendType != BackendTypeCloudStorage {
		return nil, fmt.Errorf("backend is not of type cloud_storage")
	}
	backend, ok := gc.Backend.(*CloudStorageBackendConfig) // Cast to pointer
	if !ok {
		return nil, fmt.Errorf("failed to cast backend to *CloudStorageBackendConfig")
	}
	return backend, nil
}

// PostgresBackend returns the PostgresBackendConfig from GlobalConfig
func (gc *GlobalConfig) PostgresBackend() (*PostgresBackendConfig, error) {
	if gc.BackendType != BackendTypePostgres {
		return nil, fmt.Errorf("backend is not of type postgres")
	}
	backend, ok := gc.Backend.(*PostgresBackendConfig) // Cast to pointer
	if !ok {
		return nil, fmt.Errorf("failed to cast backend to *PostgresBackendConfig")
	}
	return backend, nil
}

// LocalBackend returns the LocalBackendConfig from GlobalConfig
func (gc *GlobalConfig) LocalBackend() (*LocalBackendConfig, error) {
	if gc.BackendType != LocalBackendType {
		return nil, fmt.Errorf("backend is not of type local")
	}
	backend, ok := gc.Backend.(*LocalBackendConfig) // Cast to pointer
	if !ok {
		return nil, fmt.Errorf("failed to cast backend to *LocalBackendConfig")
	}
	return backend, nil
}
