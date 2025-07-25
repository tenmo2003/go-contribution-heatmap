package utils

import (
	"os"
	"strings"
)

func ExpandPath(path string) string {
	if strings.HasPrefix(path, ".") {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		return strings.Replace(path, ".", wd, 1)
	}
	if strings.HasPrefix(path, "~") {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		return strings.Replace(path, "~", userHomeDir, 1)
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
