package data

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	User               string `yaml:"user"`
	Pass               string `yaml:"pass"`
	Host               string `yaml:"host"`
	Port               string `yaml:"port"`
	DBName             string `yaml:"dbName"`
	MaxOpenConnections int    `yaml:"maxOpenConnections"`
	MaxIdleConnections int    `yaml:"maxIdleConnections"`
	MaxLifetime        int    `yaml:"maxLifetime"`
}

const fileName string = "config.yaml"

func NewConfigFromFile() *Config {
	var r Config
	configFile, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("opening config file %v \n", err.Error())
		return &r
	}

	defer configFile.Close()
	// jsonParser := json.NewDecoder(configFile)
	// if err = jsonParser.Decode(&r); err != nil {
	// 	fmt.Printf("parsing config file %v \n", err.Error())
	// }

	yamlParser := yaml.NewDecoder(configFile)
	if err = yamlParser.Decode(&r); err != nil {
		fmt.Printf("parsing config file %v \n", err.Error())
	}

	return &r
}
