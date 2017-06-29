package main

import (
	"fmt"
	"io/ioutil"

	"github.com/go-yaml/yaml"
	"github.com/xackery/eqemuconfig"
)

type YamlConfig struct {
	Www       string `yaml:"www,omitempty"`
	Output    string `yaml:"output,omitempty"`
	Templates string `yaml:"templates,omitempty"`
	Quests    string `yaml:"quests,omitempty"`
}

func loadYamlConfig() (config *YamlConfig, err error) {
	inFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		err = fmt.Errorf("Failed to read config.yml: %s", err.Error())
		return
	}
	config = &YamlConfig{}
	if err = yaml.Unmarshal(inFile, &config); err != nil {
		err = fmt.Errorf("Failed to unmarshal: %s", err.Error())
		return
	}
	return
}

func loadEqemuConfig() (config *eqemuconfig.Config, err error) {
	config = &eqemuconfig.Config{}

	if config, err = eqemuconfig.GetConfig(); err != nil {
		err = fmt.Errorf("Failed to load eqemuconfig: %s", err.Error())
		return
	}
	return
}
