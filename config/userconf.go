package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

func UnmarshalJson(target interface{}, homeBasedPath string) error {
	user, _ := user.Current()
	settingsPath := path.Join(user.HomeDir, homeBasedPath)
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
	user, _ := user.Current()
	settingsPath := path.Join(user.HomeDir, homeBasedPath)
	_, err := os.Stat(settingsPath)
	if err != nil {
		return NoFile(settingsPath)
	}

	settingsBin, errSettings := json.Marshal(target)
	if nil != errSettings {
		return errSettings
	}

	settingsErr := ioutil.WriteFile(settingsPath, settingsBin, 0660)
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
