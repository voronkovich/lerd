package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/geodro/lerd/internal/config"
	"github.com/geodro/lerd/internal/laravel"
	"github.com/geodro/lerd/internal/nginx"
	nodeDet "github.com/geodro/lerd/internal/node"
	phpDet "github.com/geodro/lerd/internal/php"
	"github.com/geodro/lerd/internal/podman"
	"github.com/spf13/cobra"
)

// NewParkCmd returns the park command.
func NewParkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "park [directory]",
		Short: "Park a directory to serve all subdirectories as sites",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runPark,
	}
}

func runPark(_ *cobra.Command, args []string) error {
	dir := ""
	if len(args) > 0 {
		dir = args[0]
	} else {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	// Resolve absolute path
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	cfg, err := config.LoadGlobal()
	if err != nil {
		return err
	}

	fmt.Printf("Parking directory: %s\n", absDir)

	// Add to parked directories in global config
	found := false
	for _, pd := range cfg.ParkedDirectories {
		if pd == absDir {
			found = true
			break
		}
	}
	if !found {
		cfg.ParkedDirectories = append(cfg.ParkedDirectories, absDir)
		if err := config.SaveGlobal(cfg); err != nil {
			return err
		}
	}

	// Scan subdirectories
	entries, err := os.ReadDir(absDir)
	if err != nil {
		return err
	}

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		projectDir := filepath.Join(absDir, entry.Name())
		if !laravel.IsLaravel(projectDir) {
			continue
		}

		name := entry.Name()
		domain := name + "." + cfg.DNS.TLD

		phpVersion, err := phpDet.DetectVersion(projectDir)
		if err != nil {
			phpVersion = cfg.PHP.DefaultVersion
		}

		nodeVersion, err := nodeDet.DetectVersion(projectDir)
		if err != nil {
			nodeVersion = cfg.Node.DefaultVersion
		}

		site := config.Site{
			Name:        name,
			Domain:      domain,
			Path:        projectDir,
			PHPVersion:  phpVersion,
			NodeVersion: nodeVersion,
			Secured:     false,
		}

		if err := config.AddSite(site); err != nil {
			fmt.Printf("  [WARN] could not register %s: %v\n", name, err)
			continue
		}

		if err := nginx.GenerateVhost(site, phpVersion); err != nil {
			fmt.Printf("  [WARN] could not generate vhost for %s: %v\n", name, err)
			continue
		}

		if err := ensureFPMQuadlet(phpVersion); err != nil {
			fmt.Printf("  [WARN] could not ensure FPM for PHP %s: %v\n", phpVersion, err)
		}

		fmt.Printf("  + %s -> %s (PHP %s, Node %s)\n", name, domain, phpVersion, nodeVersion)
		count++
	}

	if count > 0 {
		fmt.Printf("Reloading nginx (%d sites registered)...\n", count)
		if err := nginx.Reload(); err != nil {
			fmt.Printf("  [WARN] nginx reload: %v\n", err)
		}
	} else {
		fmt.Println("No Laravel projects found in directory.")
	}

	return nil
}

// ensureFPMQuadlet writes a PHP-FPM quadlet for the given version if it doesn't exist.
func ensureFPMQuadlet(phpVersion string) error {
	versionShort := strings.ReplaceAll(phpVersion, ".", "")
	unitName := "lerd-php" + versionShort + "-fpm"

	// Check if quadlet already exists
	quadletDir := config.QuadletDir()
	quadletPath := filepath.Join(quadletDir, unitName+".container")
	if _, err := os.Stat(quadletPath); err == nil {
		return nil // already exists
	}

	tmplContent, err := podman.GetQuadletTemplate("lerd-php-fpm.container.tmpl")
	if err != nil {
		return err
	}

	// Simple string replacement for the template
	content := strings.ReplaceAll(tmplContent, "{{.Version}}", phpVersion)
	content = strings.ReplaceAll(content, "{{.VersionShort}}", versionShort)
	content = strings.ReplaceAll(content, "{{.ProjectsDir}}", filepath.Dir(config.DataDir()))

	if err := podman.WriteQuadlet(unitName, content); err != nil {
		return err
	}

	if err := podman.DaemonReload(); err != nil {
		return err
	}

	return podman.StartUnit(unitName)
}
