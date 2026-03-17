package systemd

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed units
var unitsFS embed.FS

// GetUnit returns the content of an embedded systemd unit file.
func GetUnit(name string) (string, error) {
	// name may or may not have .service suffix
	filename := name
	if !strings.HasSuffix(filename, ".service") {
		filename += ".service"
	}
	data, err := unitsFS.ReadFile("units/" + filename)
	if err != nil {
		return "", fmt.Errorf("systemd unit %q not found: %w", name, err)
	}
	return string(data), nil
}
