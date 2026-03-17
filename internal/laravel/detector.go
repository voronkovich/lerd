package laravel

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// IsLaravel returns true if the given directory looks like a Laravel project.
func IsLaravel(dir string) bool {
	// Check for artisan file
	if _, err := os.Stat(filepath.Join(dir, "artisan")); err == nil {
		return true
	}

	// Check composer.json for laravel/framework
	composerFile := filepath.Join(dir, "composer.json")
	data, err := os.ReadFile(composerFile)
	if err == nil {
		var composer struct {
			Require    map[string]string `json:"require"`
			RequireDev map[string]string `json:"require-dev"`
		}
		if json.Unmarshal(data, &composer) == nil {
			if _, ok := composer.Require["laravel/framework"]; ok {
				return true
			}
			if _, ok := composer.RequireDev["laravel/framework"]; ok {
				return true
			}
		}
	}

	// Check public/index.php for Illuminate reference
	indexFile := filepath.Join(dir, "public", "index.php")
	indexData, err := os.ReadFile(indexFile)
	if err == nil {
		if strings.Contains(string(indexData), `Illuminate\`) || strings.Contains(string(indexData), "Illuminate\\") {
			return true
		}
	}

	return false
}
