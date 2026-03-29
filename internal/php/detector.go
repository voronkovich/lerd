package php

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/geodro/lerd/internal/config"
	"gopkg.in/yaml.v3"
)

// DetectExtensions reads composer.json in dir and returns the list of PHP extensions
// declared in the require map (ext-* keys), with the "ext-" prefix stripped.
// Returns an empty slice on any error (non-fatal).
func DetectExtensions(dir string) []string {
	data, err := os.ReadFile(filepath.Join(dir, "composer.json"))
	if err != nil {
		return nil
	}
	var composer struct {
		Require map[string]string `json:"require"`
	}
	if err := json.Unmarshal(data, &composer); err != nil {
		return nil
	}
	var exts []string
	for key := range composer.Require {
		if strings.HasPrefix(key, "ext-") {
			exts = append(exts, strings.TrimPrefix(key, "ext-"))
		}
	}
	return exts
}

// DetectVersion detects the PHP version for the given directory.
// It checks, in order:
//  1. .lerd.yaml php_version field (explicit lerd override)
//  2. .php-version file (explicit per-project pin)
//  3. composer.json require.php semver (project requirement)
//  4. global config default
func DetectVersion(dir string) (string, error) {
	// 1. .lerd.yaml — explicit lerd override takes top priority
	lerdYaml := filepath.Join(dir, ".lerd.yaml")
	if data, err := os.ReadFile(lerdYaml); err == nil {
		var lerdCfg struct {
			PHPVersion string `yaml:"php_version"`
		}
		if yaml.Unmarshal(data, &lerdCfg) == nil && lerdCfg.PHPVersion != "" {
			return bestSupportedVersion("=" + lerdCfg.PHPVersion)
		}
	}

	// 2. .php-version file — explicit per-project pin
	phpVersionFile := filepath.Join(dir, ".php-version")
	if data, err := os.ReadFile(phpVersionFile); err == nil {
		v := strings.TrimSpace(string(data))
		if v != "" {
			return bestSupportedVersion("=" + v)
		}
	}

	// 3. composer.json require.php — project requirement
	composerFile := filepath.Join(dir, "composer.json")
	if data, err := os.ReadFile(composerFile); err == nil {
		var composer struct {
			Require map[string]string `json:"require"`
		}
		if json.Unmarshal(data, &composer) == nil {
			if phpConstraint, ok := composer.Require["php"]; ok {
				return bestSupportedVersion(phpConstraint)
			}
		}
	}

	// 4. global config default
	cfg, err := config.LoadGlobal()
	if err == nil && cfg.PHP.DefaultVersion != "" {
        return bestSupportedVersion("=" + cfg.PHP.DefaultVersion)
	}

    // Return the latest supported version
    return bestSupportedVersion("*")
}

// bestSupportedVersion finds the highest supported PHP version that satisfies the given constraint.
// Returns an error if the constraint is invalid or no supported version matches.
func bestSupportedVersion(constraintStr string) (string, error) {
	constraint, err := semver.NewConstraint(constraintStr)
	if err != nil {
		return "", fmt.Errorf("invalid PHP version constraint '%s': %w", constraintStr, err)
	}

	var bestVersion *semver.Version
	for _, v := range SupportedVersions {
		version, err := semver.NewVersion(v)
		if err != nil {
			continue // skip invalid versions in list
		}
		if constraint.Check(version) {
			if bestVersion == nil || version.GreaterThan(bestVersion) {
				bestVersion = version
			}
		}
	}

	if bestVersion == nil {
        return "", fmt.Errorf("no supported PHP version satisfies constraint '%s'", constraintStr)
	}

    return bestVersion.String(), nil
}
