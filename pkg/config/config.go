package config

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Config struct {
	DefaultDestUrl string              `json:"default_dest_url"`
	Mappings       []map[string]string `json:"mappings"`
	Loggers        []map[string]string `json:"loggers,omitempty"`
}

var conf Config

func Init(jsonConfData []byte) error {
	err := json.Unmarshal(jsonConfData, &conf)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to unmarshal json conf data. err:%s", err))
	}
	return nil
}

func DefaultDestUrl() string {
	return conf.DefaultDestUrl
}

func URLMappings() []map[string]string {
	return conf.Mappings
}

func Loggers() []map[string]string {
	return conf.Loggers
}
