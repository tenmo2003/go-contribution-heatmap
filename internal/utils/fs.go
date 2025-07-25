package utils

import (
	"os"
	"strings"
)

func ExpandPath(path string) string {
	if path == "." {
		return os.Getenv("PWD")
	}
	if strings.HasPrefix(path, "~") {
		return strings.Replace(path, "~", os.Getenv("HOME"), 1)
	}
	return path
}

func ExpandPaths(paths []string) []string {
	expandedPaths := []string{}
	for _, path := range paths {
		expandedPaths = append(expandedPaths, ExpandPath(path))
	}
	return expandedPaths
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
