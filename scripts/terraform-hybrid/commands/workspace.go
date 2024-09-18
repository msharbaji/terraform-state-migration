package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const providerBaseDir = "deploy/provider"

// WorkspaceCmd defines the structure for the Workspace command
type WorkspaceCmd struct {
	SelectOrCreate bool   `help:"Select a workspace, or create it if it doesn't exist."`
	Select         string `help:"Select an existing workspace."`
	New            string `help:"Create a new workspace."`
	List           bool   `help:"List all available workspaces."`
	Current        bool   `help:"Show the current active workspace."`
	Delete         string `help:"Delete a workspace."`
}

// Run executes the logic for the Workspace command
func (w *WorkspaceCmd) Run() error {
	switch {
	case w.List:
		return w.ListWorkspaces()
	case w.Current:
		return w.CurrentWorkspace()
	case w.New != "":
		return w.CreateWorkspace(w.New)
	case w.Select != "":
		return w.SelectWorkspace(w.Select)
	case w.Delete != "":
		return w.DeleteWorkspace(w.Delete)
	case w.SelectOrCreate:
		return w.SelectOrCreateWorkspace()
	default:
		return fmt.Errorf("no workspace operation provided, use --help for options")
	}
}

// ListWorkspaces lists all available workspaces
func (w *WorkspaceCmd) ListWorkspaces() error {
	return runTerraformCommand("Listing available workspaces...", "terraform", "workspace", "list")
}

// CurrentWorkspace shows the current workspace
func (w *WorkspaceCmd) CurrentWorkspace() error {
	return runTerraformCommand("Showing current workspace...", "terraform", "workspace", "show")
}

// CreateWorkspace creates a new workspace
func (w *WorkspaceCmd) CreateWorkspace(workspace string) error {
	return runTerraformCommand(
		fmt.Sprintf("Creating new workspace: %s", workspace),
		"terraform", "workspace", "new", workspace)
}

// SelectWorkspace selects the specified workspace
func (w *WorkspaceCmd) SelectWorkspace(workspace string) error {
	return runTerraformCommand(
		fmt.Sprintf("Selecting workspace: %s", workspace),
		"terraform", "workspace", "select", workspace)
}

// SelectOrCreateWorkspace selects or creates a workspace based on the current directory
func (w *WorkspaceCmd) SelectOrCreateWorkspace() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current directory: %v", err)
	}

	baseDirIndex := strings.Index(currentDir, providerBaseDir)
	if baseDirIndex == -1 {
		return fmt.Errorf("current directory is not within the %s directory", providerBaseDir)
	}

	baseDir := currentDir[:baseDirIndex+len(providerBaseDir)]
	relativePath, err := filepath.Rel(baseDir, currentDir)
	if err != nil {
		return fmt.Errorf("error calculating relative path: %v", err)
	}

	workspace := strings.ReplaceAll(relativePath, string(filepath.Separator), "_")
	workspace = strings.TrimPrefix(workspace, string(filepath.Separator))

	fmt.Printf("Selecting or creating workspace: %s\n", workspace)

	return runTerraformCommand(
		fmt.Sprintf("Selecting workspace: %s", workspace),
		"terraform", "workspace", "select", "--or-create",
		workspace)

}

// DeleteWorkspace deletes the specified workspace
func (w *WorkspaceCmd) DeleteWorkspace(workspace string) error {
	return runTerraformCommand(
		fmt.Sprintf("Deleting workspace: %s", workspace),
		"terraform", "workspace", "delete", workspace)
}

func runTerraformCommand(message string, args ...string) error {
	fmt.Println(message)
	fmt.Printf("Running terraform command: %s\n", strings.Join(args, " "))
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		return fmt.Errorf("error executing command: %v, output: %s", err, string(output))
	}
	return nil
}
