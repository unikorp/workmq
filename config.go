package workmq

import (
	"encoding/json"
	"fmt"
	"os"
)

// QueueConfig is the "queues" configuration section type definition.
type QueueConfig struct {
	Processor  string `json:"processor"`
	NumWorkers int    `json:"num_workers"`
}

// PortsConfig is the "port" configuration section type definition.
type PortsConfig struct {
	UDP  string `json:"udp"`
	HTTP string `json:"http"`
}

// Config is the configuration type definition.
type Config struct {
	Ports  PortsConfig            `json:"ports"`
	Queues map[string]QueueConfig `json:"queues"`
}

// GetConfig returns the configuration object that can be used anywhere in application.
func GetConfig() Config {
	file, _ := os.Open("./config.json")
	decoder := json.NewDecoder(file)

	config := Config{}
	err := decoder.Decode(&config)

	if err != nil {
		fmt.Println("An error occurs on configuration loading:", err)
	}

	return config
}
