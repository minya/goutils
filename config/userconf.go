// Package config provides utilities for handling user configuration files.
// It includes functions for marshaling and unmarshaling JSON configuration files,
// with support for both relative and home-based (~/) paths.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// UnmarshalJson reads a JSON file from the specified path and deserializes it into the target object.
// The path can be either relative to the current directory or home-based (starting with "~/").
// Returns a NoFileError (which matches fs.ErrNotExist) only when the file does not exist; any
// other failure, such as permission denied, is returned as is, so callers do not mistake an
// unreadable file for a missing one.
func UnmarshalJson(target interface{}, relativeOrHomeBasedPath string) error {
	settingsPath := buildPath(relativeOrHomeBasedPath)
	settingsBin, err := os.ReadFile(settingsPath)
	if errors.Is(err, fs.ErrNotExist) {
		return &NoFileError{msg: "No file exists: " + settingsPath, err: err}
	}
	if err != nil {
		return fmt.Errorf("reading %s: %w", settingsPath, err)
	}
	if err := json.Unmarshal(settingsBin, target); err != nil {
		return fmt.Errorf("parsing %s: %w", settingsPath, err)
	}
	return nil
}

// MarshalJson serializes the target object to JSON and writes it to the specified file path.
// The path can be either relative to the current directory or home-based (starting with "~/").
// Settings usually hold tokens and passwords, so the file is written with mode 0600, and it is
// replaced atomically: the JSON goes to a temporary file in the same directory, which is synced
// and then renamed over the target, so a crash never leaves a truncated file.
func MarshalJson(target interface{}, relativeOrHomeBasedPath string) error {
	settingsPath := buildPath(relativeOrHomeBasedPath)
	settingsBin, err := json.Marshal(target)
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(settingsPath), "."+filepath.Base(settingsPath)+".tmp-*")
	if err != nil {
		return fmt.Errorf("unable to create %v: %w", settingsPath, err)
	}
	defer os.Remove(tmp.Name()) // no-op once renamed
	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return fmt.Errorf("writing %v: %w", settingsPath, err)
	}
	if _, err := tmp.Write(settingsBin); err != nil {
		tmp.Close()
		return fmt.Errorf("writing %v: %w", settingsPath, err)
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return fmt.Errorf("writing %v: %w", settingsPath, err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("writing %v: %w", settingsPath, err)
	}
	if err := os.Rename(tmp.Name(), settingsPath); err != nil {
		return fmt.Errorf("writing %v: %w", settingsPath, err)
	}
	return nil
}

// buildPath converts a path that might be relative or home-based (starting with "~/")
// to an absolute path. Home-based paths are expanded using the HOME environment variable.
func buildPath(relativeOrHomeBasedPath string) string {
	if strings.HasPrefix(relativeOrHomeBasedPath, "~/") {
		home := os.Getenv("HOME")
		trimmedHomeBased := strings.TrimPrefix(relativeOrHomeBasedPath, "~/")
		return path.Join(home, trimmedHomeBased)
	}
	return relativeOrHomeBasedPath
}

// NoFileError is returned when an operation is attempted on a file that doesn't exist.
type NoFileError struct {
	msg string
	err error
}

// Error implements the error interface for NoFileError.
func (e *NoFileError) Error() string {
	return e.msg
}

// Unwrap returns the underlying error; it matches fs.ErrNotExist.
func (e *NoFileError) Unwrap() error {
	if e.err == nil {
		return fs.ErrNotExist
	}
	return e.err
}

// NoFile creates a new NoFileError for the specified filename.
func NoFile(fileName string) error {
	return &NoFileError{msg: "No file exists: " + fileName, err: fs.ErrNotExist}
}
