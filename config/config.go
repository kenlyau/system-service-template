package config

import (
	"encoding/json"
	"io/ioutil"
	"system-service-template/utils"
)

var _config Config

type Config struct {
	ServiceName        string
	ServiceDisplayName string
	ServiceDescription string
	HttpPort           string
	LogName            string
	LogMaxSize         int
	LogMaxAge          int
}

func LoadConfig() error {
	configPath := utils.GetAbsPath("config.json")
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(b, &_config); err != nil {
		return err
	}
	return nil
}

func GetConfig() Config {
	return _config
}
