package config

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "log"
)

type Config struct {
    Source string `yaml:"source"`
    Extensions []string `yaml:"extensions"`
    Dest string `yaml:"dest"`
    MinSize int64 `yaml:"min_size"`
    Aliases map[string]string `yaml:"aliases"`
}



func (c *Config) GetConfig() *Config {
    yamlFile, err := ioutil.ReadFile("config.yaml")
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal(yamlFile, c)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }

    return c
}
