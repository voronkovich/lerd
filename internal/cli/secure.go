package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/geodro/lerd/internal/certs"
	"github.com/geodro/lerd/internal/config"
	"github.com/geodro/lerd/internal/nginx"
	"github.com/spf13/cobra"
)

// NewSecureCmd returns the secure command.
func NewSecureCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "secure [name]",
		Short: "Enable HTTPS for the current site using mkcert",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runSecure,
	}
}

func runSecure(_ *cobra.Command, args []string) error {
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
		return fmt.Errorf("site %q not found — run 'lerd link' first", name)
	}

	certsDir := filepath.Join(config.CertsDir(), "sites")
	fmt.Printf("Issuing certificate for %s...\n", site.Domain)

	if err := certs.IssueCert(site.Domain, certsDir); err != nil {
		return fmt.Errorf("issuing certificate: %w", err)
	}

	// Generate SSL vhost
	if err := nginx.GenerateSSLVhost(*site, site.PHPVersion); err != nil {
		return fmt.Errorf("generating SSL vhost: %w", err)
	}

	// Remove plain HTTP vhost (replaced by redirect in SSL template)
	if err := nginx.RemoveVhost(site.Domain); err != nil {
		fmt.Printf("[WARN] removing HTTP vhost: %v\n", err)
	}

	// Rename SSL vhost to replace HTTP one
	sslConf := filepath.Join(config.NginxConfD(), site.Domain+"-ssl.conf")
	mainConf := filepath.Join(config.NginxConfD(), site.Domain+".conf")
	if err := os.Rename(sslConf, mainConf); err != nil {
		return fmt.Errorf("renaming SSL config: %w", err)
	}

	// Mark site as secured
	site.Secured = true
	if err := config.AddSite(*site); err != nil {
		return fmt.Errorf("updating site registry: %w", err)
	}

	fmt.Printf("Secured: https://%s\n", site.Domain)

	if err := nginx.Reload(); err != nil {
		fmt.Printf("[WARN] nginx reload: %v\n", err)
	}

	return nil
}
