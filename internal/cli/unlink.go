package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/geodro/lerd/internal/config"
	"github.com/geodro/lerd/internal/nginx"
	"github.com/spf13/cobra"
)

// NewUnlinkCmd returns the unlink command.
func NewUnlinkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unlink [name]",
		Short: "Unlink the current directory site",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runUnlink,
	}
}

func runUnlink(_ *cobra.Command, args []string) error {
	name := ""
	if len(args) > 0 {
		name = args[0]
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		name = filepath.Base(cwd)
	}

	site, err := config.FindSite(name)
	if err != nil {
		return fmt.Errorf("site %q not found", name)
	}

	if err := nginx.RemoveVhost(site.Domain); err != nil {
		fmt.Printf("[WARN] removing vhost: %v\n", err)
	}

	if err := config.RemoveSite(name); err != nil {
		return fmt.Errorf("removing site from registry: %w", err)
	}

	fmt.Printf("Unlinked: %s (%s)\n", name, site.Domain)

	if err := nginx.Reload(); err != nil {
		fmt.Printf("[WARN] nginx reload: %v\n", err)
	}

	return nil
}
