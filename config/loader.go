package config

import (
	"encoding/json"
	"log"
	"os"
)

var config *AppConfig = nil

// LoadConfig loads the config file if needed, or returns the active config if we have it in memory.
func LoadConfig() AppConfig {
	if config != nil {
		return *config
	}

	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatalln("A config file is required.")
	}
	data := new([]byte)
	_, err = file.Read(*data)
	if err != nil {
		log.Fatalf("An error occured while reading config file: %s\n", err)
	}

	err = json.Unmarshal(*data, config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %s\n", err)
	}

	return *config
}
