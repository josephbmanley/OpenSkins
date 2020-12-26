package pluginmanager

import (
	"os"
	"path/filepath"
)

// GetPlugins gets plugin files in directory
func GetPlugins(pluginDir string) ([]string, error) {
	pattern := "*.so"
	var matches []string
	err := filepath.Walk(pluginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
