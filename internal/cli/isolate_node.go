package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/geodro/lerd/internal/config"
	"github.com/spf13/cobra"
)

// NewIsolateNodeCmd returns the isolate:node command.
func NewIsolateNodeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "isolate:node <version>",
		Short: "Pin the Node.js version for the current directory",
		Args:  cobra.ExactArgs(1),
		RunE:  runIsolateNode,
	}
}

func runIsolateNode(_ *cobra.Command, args []string) error {
	version := args[0]

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	nodeVersionFile := filepath.Join(cwd, ".node-version")
	if err := os.WriteFile(nodeVersionFile, []byte(version+"\n"), 0644); err != nil {
		return fmt.Errorf("writing .node-version: %w", err)
	}

	fmt.Printf("Node.js version pinned to %s in %s\n", version, cwd)

	// Run fnm install for this version
	fnmPath := filepath.Join(config.BinDir(), "fnm")
	if _, err := os.Stat(fnmPath); err == nil {
		cmd := exec.Command(fnmPath, "install", version)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("[WARN] fnm install %s: %v\n", version, err)
		}
	} else {
		fmt.Println("[WARN] fnm not found — run 'lerd install' to set up Node.js management")
	}

	return nil
}
