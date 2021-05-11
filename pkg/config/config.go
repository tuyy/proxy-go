package config

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Config struct {
	defaultDestUrl string              `json:"default_dest_url"`
	mappings       []map[string]string `json:"mappings"`
	loggers        []map[string]string `json:"loggers,omitempty"`
}

var conf Config

func DefaultDestUrl() string {
	return conf.defaultDestUrl
}

func URLMappings() []map[string]string {
	return conf.mappings
}

func Loggers() []map[string]string {
	return conf.loggers
}

func Init(jsonConfData []byte) error {
	err := json.Unmarshal(jsonConfData, &conf)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to unmarshal json conf data. err:%s", err))
	}
	return nil
}
