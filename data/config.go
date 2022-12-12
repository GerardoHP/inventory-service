package data

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Host   string `json:"host"`
	Port   string `json:"port"`
	DBName string `json:"dbName"`
}

const fileName string = "config.json"

func NewConfigFromFile() *Config {
	var r Config
	configFile, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("opening config file %v \n", err.Error())
		return &r
	}

	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&r); err != nil {
		fmt.Printf("parsing config file %v \n", err.Error())
	}

	return &r
}
