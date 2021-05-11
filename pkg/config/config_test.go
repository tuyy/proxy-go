package config

import (
	"encoding/json"
	"fmt"
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
	]
}`)

func TestUnmarshal(t *testing.T) {
	var jsonBlob = []byte(`[     {"Name": "Platypus", "Order": "Monotremata"},     {"Name": "Quoll",    "Order": "Dasyuromorphia"} ]`)
	type Animal struct {
		Name  string
		Order string
	}
	var animals []Animal
	err := json.Unmarshal(jsonBlob, &animals)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", animals)
}

func TestInitConfigHappy(t *testing.T) {
	err := Init([]byte(happyJsonConfData))
	assert.Nil(t, err)

	expected := Config{
		defaultDestUrl: "dest1.com",
		mappings: []map[string]string{
			{"path": "/test1", "dest_url": "dest2.com"},
			{"path": "/test2", "dest_url": "dest1.com"},
			{"path": "/test3", "dest_url": "dest2.com"},
		},
	}
	assert.Equal(t, expected, conf)
}
