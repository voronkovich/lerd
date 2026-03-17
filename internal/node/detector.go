package node

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/geodro/lerd/internal/config"
)

// DetectVersion detects the Node.js version for the given directory.
// It checks, in order:
//  1. .nvmrc
//  2. .node-version
//  3. package.json engines.node
//  4. global config default
func DetectVersion(dir string) (string, error) {
	// 1. .nvmrc
	nvmrc := filepath.Join(dir, ".nvmrc")
	if data, err := os.ReadFile(nvmrc); err == nil {
		v := strings.TrimSpace(string(data))
		// strip leading 'v'
		v = strings.TrimPrefix(v, "v")
		if v != "" {
			return extractMajor(v), nil
		}
	}

	// 2. .node-version
	nodeVersion := filepath.Join(dir, ".node-version")
	if data, err := os.ReadFile(nodeVersion); err == nil {
		v := strings.TrimSpace(string(data))
		v = strings.TrimPrefix(v, "v")
		if v != "" {
			return extractMajor(v), nil
		}
	}

	// 3. package.json engines.node
	pkgJSON := filepath.Join(dir, "package.json")
	if data, err := os.ReadFile(pkgJSON); err == nil {
		var pkg struct {
			Engines struct {
				Node string `json:"node"`
			} `json:"engines"`
		}
		if json.Unmarshal(data, &pkg) == nil && pkg.Engines.Node != "" {
			if v := parseNodeConstraint(pkg.Engines.Node); v != "" {
				return v, nil
			}
		}
	}

	// 4. global config default
	cfg, err := config.LoadGlobal()
	if err != nil {
		return "22", nil
	}
	return cfg.Node.DefaultVersion, nil
}

// extractMajor returns the major version number from a semver-like string.
// e.g. "18.12.0" → "18", "22" → "22"
func extractMajor(v string) string {
	parts := strings.SplitN(v, ".", 2)
	return parts[0]
}

// parseNodeConstraint extracts the first numeric major version from a constraint.
func parseNodeConstraint(constraint string) string {
	re := regexp.MustCompile(`(\d+)`)
	m := re.FindStringSubmatch(constraint)
	if len(m) > 1 {
		return m[1]
	}
	return ""
}
