package watcher

import (
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Watch monitors the given directories for new project subdirectories.
// When an artisan file appears in a new subdirectory, onNew is called with
// the project path.
func Watch(dirs []string, onNew func(path string)) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()

	for _, dir := range dirs {
		expanded := expandHome(dir)
		if err := os.MkdirAll(expanded, 0755); err != nil {
			continue
		}
		if err := w.Add(expanded); err != nil {
			continue
		}
	}

	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				return nil
			}
			if event.Op&(fsnotify.Create|fsnotify.Write) != 0 {
				if filepath.Base(event.Name) == "artisan" {
					projectDir := filepath.Dir(event.Name)
					onNew(projectDir)
				}
			}
		case err, ok := <-w.Errors:
			if !ok {
				return nil
			}
			_ = err // log if needed
		}
	}
}

func expandHome(path string) string {
	if len(path) > 1 && path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}
