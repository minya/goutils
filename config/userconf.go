// Package config provides utilities for handling user configuration files.
// It includes functions for marshaling and unmarshaling JSON configuration files,
// with support for both relative and home-based (~/) paths.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

// UnmarshalJson reads a JSON file from the specified path and deserializes it into the target object.
// The path can be either relative to the current directory or home-based (starting with "~/").
// Returns a NoFileError if the file doesn't exist or another error if reading or unmarshaling fails.
func UnmarshalJson(target interface{}, relativeOrHomeBasedPath string) error {
	settingsPath := buildPath(relativeOrHomeBasedPath)
	_, err := os.Stat(settingsPath)
	if err != nil {
		return NoFile(settingsPath)
	}

	settingsBin, settingsErr := os.ReadFile(settingsPath)
	if settingsErr != nil {
		return settingsErr
	}
	errSettings := json.Unmarshal(settingsBin, target)
	if nil != errSettings {
		return errSettings
	}
	return nil
}

// MarshalJson serializes the target object to JSON and writes it to the specified file path.
// The path can be either relative to the current directory or home-based (starting with "~/").
// If the file doesn't exist, it will be created.
// Returns an error if the file creation, marshaling, or writing operations fail.
func MarshalJson(target interface{}, relativeOrHomeBasedPath string) error {
	settingsPath := buildPath(relativeOrHomeBasedPath)
	_, err := os.Stat(settingsPath)
	if err != nil {
		f, err := os.Create(settingsPath)
		if nil != err {
			return fmt.Errorf("unable to create %v", settingsPath)
		} else {
			f.Close()
		}
	}

	settingsBin, errSettings := json.Marshal(target)
	if nil != errSettings {
		return errSettings
	}

	settingsErr := os.WriteFile(settingsPath, settingsBin, 0644)
	if settingsErr != nil {
		return settingsErr
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
}

// Error implements the error interface for NoFileError.
func (e *NoFileError) Error() string {
	return e.msg
}

// NoFile creates a new NoFileError for the specified filename.
func NoFile(fileName string) error {
	return &NoFileError{"No file exists: " + fileName}
}
