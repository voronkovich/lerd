package php

import (
	"path/filepath"
	"regexp"
	"sort"

	"github.com/geodro/lerd/internal/config"
)

// SupportedVersions lists the PHP versions that lerd supports.
var SupportedVersions = []string{"8.1", "8.2", "8.3", "8.4", "8.5"}

var fpmQuadletRe = regexp.MustCompile(`^lerd-php(\d)(\d+)-fpm\.container$`)

// ListInstalled returns all PHP versions that have an FPM quadlet file,
// e.g. ["8.3", "8.4"]. These are installed via `lerd use <version>`.
func ListInstalled() ([]string, error) {
	pattern := filepath.Join(config.QuadletDir(), "lerd-php*-fpm.container")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var versions []string
	for _, m := range matches {
		base := filepath.Base(m)
		sub := fpmQuadletRe.FindStringSubmatch(base)
		if len(sub) == 3 {
			versions = append(versions, sub[1]+"."+sub[2])
		}
	}

	sort.Strings(versions)
	return versions, nil
}

// IsInstalled returns true if the given PHP version has an FPM quadlet.
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
