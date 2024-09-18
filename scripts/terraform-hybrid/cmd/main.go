package main

import (
	"log"

	"github.com/msharbaji/terraform-state-migration/terraform-hybrid/commands"

	"github.com/alecthomas/kong"
)

// CLI represents the structure for CLI commands
//
//nolint:lll
var CLI struct {
	GenerateBackend commands.GenerateBackendCmd `cmd:"" help:"Generate backend.tf files for a given config and provider folder."`
	Workspace       commands.WorkspaceCmd       `cmd:"" help:"Manage Terraform workspaces (create, select, list, delete)."`
}

func main() {
	// Initialize Kong and parse CLI arguments
	ctx := kong.Parse(&CLI)

	// Run the appropriate command handler
	err := ctx.Run(&CLI)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
