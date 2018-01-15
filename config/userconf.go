package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

//UnmarshalJson reads and deserializes object from file
func UnmarshalJson(target interface{}, relativeOrHomeBasedPath string) error {
	settingsPath := buildPath(relativeOrHomeBasedPath)
	_, err := os.Stat(settingsPath)
	if err != nil {
		return NoFile(settingsPath)
	}

	settingsBin, settingsErr := ioutil.ReadFile(settingsPath)
	if settingsErr != nil {
		return settingsErr
	}
	errSettings := json.Unmarshal(settingsBin, target)
	if nil != errSettings {
		return errSettings
	}
	return nil
}

//MarshalJson serializes and writes object to file
func MarshalJson(target interface{}, relativeOrHomeBasedPath string) error {
	settingsPath := buildPath(relativeOrHomeBasedPath)
	_, err := os.Stat(settingsPath)
	if err != nil {
		f, err := os.Create(settingsPath)
		if nil != err {
			return fmt.Errorf("Unable to create %v", settingsPath)
		} else {
			f.Close()
		}
	}

	settingsBin, errSettings := json.Marshal(target)
	if nil != errSettings {
		return errSettings
	}

	settingsErr := ioutil.WriteFile(settingsPath, settingsBin, 0644)
	if settingsErr != nil {
		return settingsErr
	}

	return nil
}

func buildPath(relativeOrHomeBasedPath string) string {
	if strings.Index(relativeOrHomeBasedPath, "~/") == 0 {
		home := os.Getenv("HOME")
		trimmedHomeBased := strings.TrimLeft(relativeOrHomeBasedPath, "~/")
		return path.Join(home, trimmedHomeBased)
	}
	return relativeOrHomeBasedPath
}

type NoFileError struct {
	msg string
}

func (e *NoFileError) Error() string {
	return e.msg
}

func NoFile(fileName string) error {
	return &NoFileError{"No file exits: " + fileName}
}
