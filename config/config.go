package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// Config mirrors the config.yaml file
// See config.yaml.dist for a field by field description
type Config struct {
	Source     string            `yaml:"source"`
	Extensions []string          `yaml:"extensions"`
	Dest       string            `yaml:"dest"`
	MinSize    int64             `yaml:"min_size"`
	Aliases    map[string]string `yaml:"aliases"`
}

// GetConfig loads config.yaml into the Config data structure
func (c *Config) GetConfig(path string) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err %v\n", err)
        return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
        log.Printf("Invalid config file! %s\nError: %v\n", path, err)
        return err
	}

	return nil
}
