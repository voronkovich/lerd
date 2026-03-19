//go:build nogui

package tray

import "fmt"

// Run returns an error when the binary was built without CGo/AppIndicator support.
func Run(_ bool) error {
	return fmt.Errorf("system tray not available (built without CGo/AppIndicator support)")
}
