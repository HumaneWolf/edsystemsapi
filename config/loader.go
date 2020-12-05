package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var configLoaded = false
var config = AppConfig{}

// LoadConfig loads the config file if needed, or returns the active config if we have it in memory.
func LoadConfig() AppConfig {
	if configLoaded {
		return config
	}

	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatalln("A config file is required.")
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("An error occured while reading config file: %s\n", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %s\n", err)
	}

	configLoaded = true
	return config
}
