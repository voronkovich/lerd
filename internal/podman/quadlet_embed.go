package podman

import (
	"embed"
	"fmt"
	"strings"
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

// ApplyExtraPorts appends extra PublishPort lines to quadlet content.
func ApplyExtraPorts(content string, extraPorts []string) string {
	var sb strings.Builder
	sb.WriteString(content)
	for _, p := range extraPorts {
		fmt.Fprintf(&sb, "PublishPort=%s\n", p)
	}
	return sb.String()
}
