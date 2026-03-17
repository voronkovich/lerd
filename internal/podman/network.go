package podman

import (
	"strings"
)

// EnsureNetwork creates the named Podman network if it does not already exist.
func EnsureNetwork(name string) error {
	out, err := Run("network", "ls", "--format={{.Name}}")
	if err != nil {
		return err
	}

	for _, line := range strings.Split(out, "\n") {
		if strings.TrimSpace(line) == name {
			return nil // already exists
		}
	}

	return RunSilent("network", "create", "--driver", "bridge", name)
}
