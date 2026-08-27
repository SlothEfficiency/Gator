package config

import (
	"encoding/json"
	"os"
	"path"
)

func getConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	return path.Join(homedir, conifFileName), err
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	jsonfile, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	var configfile Config
	err = json.Unmarshal(jsonfile, &configfile)

	return configfile, err
}

func write(c Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	jsonfile, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, jsonfile, os.ModePerm)
	return err
}

func SetUser(user string) error {
	c, err := Read()
	if err != nil {
		return err
	}

	c.CurrentUserName = user
	return write(c)
}
