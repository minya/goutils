package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func UnmarshalJson(target interface{}, homeBasedPath string) error {
	settingsPath := path.Join(os.Getenv("HOME"), homeBasedPath)
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

func MarshalJson(target interface{}, homeBasedPath string) error {
	settingsPath := path.Join(os.Getenv("HOME"), homeBasedPath)
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

type NoFileError struct {
	msg string
}

func (e *NoFileError) Error() string {
	return e.msg
}

func NoFile(fileName string) error {
	return &NoFileError{"No file exits: " + fileName}
}
