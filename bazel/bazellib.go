package bazel

import (
	"os"
	"path/filepath"
)

// Workspace returns the root of the workspace when running by Bazel.
func Workspace() string {
	if os.Getenv("TEST_SRCDIR") != "" {
		return filepath.Join(
			os.Args[0]+".runfiles",
			"__main__")
	}

	return filepath.Join(
		os.Getenv("TEST_SRCDIR"),
		os.Getenv("TEST_WORKSPACE"),
	)
}
