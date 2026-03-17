package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// NewIsolateCmd returns the isolate command.
func NewIsolateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "isolate <version>",
		Short: "Pin the PHP version for the current directory",
		Args:  cobra.ExactArgs(1),
		RunE:  runIsolate,
	}
}

func runIsolate(_ *cobra.Command, args []string) error {
	version := args[0]

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	phpVersionFile := filepath.Join(cwd, ".php-version")
	if err := os.WriteFile(phpVersionFile, []byte(version+"\n"), 0644); err != nil {
		return fmt.Errorf("writing .php-version: %w", err)
	}

	fmt.Printf("PHP version pinned to %s in %s\n", version, cwd)
	return nil
}
