package utils

import (
	"os"
	"path/filepath"
)

// FolderFinder defines the interface for finding component provider folders
type FolderFinder interface {
	FindComponentProviderFolders(providerFolderPath string) ([]string, error)
}

// TerraformFolderFinder is the concrete implementation for finding folders
type TerraformFolderFinder struct{}

// NewFolderFinder creates a new instance of FolderFinder
func NewFolderFinder() FolderFinder {
	return &TerraformFolderFinder{}
}

// FindComponentProviderFolders finds all subfolders within a specified provider folder
func (tff *TerraformFolderFinder) FindComponentProviderFolders(providerFolderPath string) ([]string, error) {
	var componentFolders []string

	// Walk through the provider folder path and find all subdirectories
	err := filepath.Walk(providerFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if it's a component folder or deeper (e.g., has Terraform files)
		if info.IsDir() && filepath.Base(path) == "component" {
			componentFolders = append(componentFolders, path)
		}

		return nil
	})

	return componentFolders, err
}
