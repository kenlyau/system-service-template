package utils

import (
	"os"
	"path/filepath"
)

func GetCurrentPath() string {
	ex, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(ex)
}

func GetAbsPath(p string) string {
	return filepath.Join(GetCurrentPath(), p)
}
