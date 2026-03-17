package systemd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/geodro/lerd/internal/config"
)

// WriteService writes a systemd user service unit file.
func WriteService(name, content string) error {
	dir := config.SystemdUserDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, name+".service")
	return os.WriteFile(path, []byte(content), 0644)
}

// EnableService enables a systemd user service.
func EnableService(name string) error {
	cmd := exec.Command("systemctl", "--user", "enable", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("enable %s: %w\n%s", name, err, out)
	}
	return nil
}

// StartService starts a systemd user service.
func StartService(name string) error {
	cmd := exec.Command("systemctl", "--user", "start", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("start %s: %w\n%s", name, err, out)
	}
	return nil
}
