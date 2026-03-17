package podman

import (
	"embed"
	"fmt"
)

//go:embed quadlets
var quadletFS embed.FS

// GetQuadletTemplate returns the content of a named quadlet template file.
func GetQuadletTemplate(name string) (string, error) {
	data, err := quadletFS.ReadFile("quadlets/" + name)
	if err != nil {
		return "", fmt.Errorf("quadlet template %q not found: %w", name, err)
	}
	return string(data), nil
}
