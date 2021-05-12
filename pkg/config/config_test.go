package config

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

var happyJsonConfData = []byte(`
{
	"default_dest_url": "dest1.com",
	"mappings": [
		{"path": "/test1", "dest_url": "dest2.com"},
		{"path": "/test2", "dest_url": "dest1.com"},
		{"path": "/test3", "dest_url": "dest2.com"}
	],
	"loggers": [
        {"level": "1", "description": "debug", "filepath": "./logs/debug.log", "format": "[%D] %m", "class": "AsyncLogger"} ,
        {"level": "2", "description": "info", "filepath": "./logs/info.log", "format": "[%D] %m", "class": "AsyncLogger"}
	]
}`)

func TestInitConfigHappy(t *testing.T) {
	err := Init(happyJsonConfData)
	assert.Nil(t, err)

	expected := Config{
		DefaultDestUrl: "dest1.com",
		Mappings: []map[string]string{
			{"path": "/test1", "dest_url": "dest2.com"},
			{"path": "/test2", "dest_url": "dest1.com"},
			{"path": "/test3", "dest_url": "dest2.com"},
		},
		Loggers: []map[string]string{
			{"level": "1", "description": "debug", "filepath": "./logs/debug.log", "format": "[%D] %m", "class": "AsyncLogger"},
			{"level": "2", "description": "info", "filepath": "./logs/info.log", "format": "[%D] %m", "class": "AsyncLogger"},
		},
	}
	assert.Equal(t, expected, conf)
}

func TestStandardUnmarshal(t *testing.T) {
	var jsonBlob = []byte(`{"my_name": "tuyy", "my_order": "pizza2"}`)
	type Person struct {
		Name  string `json:"my_name"`
		Order string `json:"my_order"`
	}

	var p Person
	err := json.Unmarshal(jsonBlob, &p)
	assert.Nil(t, err)

	assert.Equal(t, "tuyy", p.Name)
	assert.Equal(t, "pizza2", p.Order)
}
