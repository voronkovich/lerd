package podman

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Run executes podman with the given arguments and returns stdout.
func Run(args ...string) (string, error) {
	cmd := exec.Command("podman", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("podman %s: %w\n%s", strings.Join(args, " "), err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

// RunSilent executes podman with the given arguments, discarding output.
func RunSilent(args ...string) error {
	_, err := Run(args...)
	return err
}

// ContainerRunning returns true if the named container is running.
func ContainerRunning(name string) (bool, error) {
	out, err := Run("inspect", "--format={{.State.Running}}", name)
	if err != nil {
		// container doesn't exist
		return false, nil
	}
	return strings.TrimSpace(out) == "true", nil
}

// ContainerExists returns true if the named container exists (running or not).
func ContainerExists(name string) (bool, error) {
	_, err := Run("inspect", "--format={{.Name}}", name)
	if err != nil {
		return false, nil
	}
	return true, nil
}
