package nginx

import (
	"embed"
	"fmt"
)

//go:embed all:templates
var templateFS embed.FS

// GetTemplate returns the content of a named nginx template/config file.
// name can be "nginx.conf", "vhost.conf.tmpl", or "vhost-ssl.conf.tmpl".
func GetTemplate(name string) ([]byte, error) {
	data, err := templateFS.ReadFile("templates/" + name)
	if err != nil {
		return nil, fmt.Errorf("nginx template %q not found: %w", name, err)
	}
	return data, nil
}
