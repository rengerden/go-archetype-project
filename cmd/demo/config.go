package main

import (
	"os"
	"encoding/json"
)


type Config struct {
	CacheTTL         int
	Providers        []string
	LimitRPP         int
	PeriodMs         int
	LimitConcurrency int
}

func GetConfig(configPath string) (cfg Config, err error) {
	configFile, err := os.Open(configPath)
	if err != nil {
		return
	}
	err = json.NewDecoder(configFile).Decode(&cfg)
	return
}
