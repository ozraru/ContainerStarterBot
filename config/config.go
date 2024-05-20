package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigStruct struct {
	Token                string `yaml:"token"`
	AppId                string `yaml:"appid"`
	Guild                string `yaml:"guild"`
	Channel              string `yaml:"channel"`
	Container            string `yaml:"container"`
	MaxTail              int64  `yaml:"max_tail"`
	Timeout              int    `yaml:"timeout"`
	EnableStop           bool   `yaml:"enable_stop"`
	EnableAutoShutdown   bool   `yaml:"enable_auto-shutdown"`
	AutoShutdownDuration string `yaml:"auto-shutdown_duration"`
}

var Config ConfigStruct

func LoadConfig() {
	buf, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("Failed to load config.yaml: ", err)
	}

	Config = ConfigStruct{}
	err = yaml.Unmarshal(buf, &Config)
	if err != nil {
		log.Fatal("Failed to unmarshal config: ", err)
	}
}
