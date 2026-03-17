package php

import (
	"path/filepath"
	"regexp"
	"sort"

	"github.com/geodro/lerd/internal/config"
)

var phpBinRe = regexp.MustCompile(`^php(\d+\.\d+)$`)

// ListInstalled returns all PHP versions installed in the lerd bin directory.
// Versions are returned in sorted order, e.g. ["8.3", "8.4"].
func ListInstalled() ([]string, error) {
	pattern := filepath.Join(config.BinDir(), "php*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var versions []string
	for _, m := range matches {
		base := filepath.Base(m)
		sub := phpBinRe.FindStringSubmatch(base)
		if len(sub) == 2 {
			versions = append(versions, sub[1])
		}
	}

	sort.Strings(versions)
	return versions, nil
}

// IsInstalled returns true if the given PHP version is available in the bin dir.
func IsInstalled(version string) bool {
	versions, err := ListInstalled()
	if err != nil {
		return false
	}
	for _, v := range versions {
		if v == version {
			return true
		}
	}
	return false
}
