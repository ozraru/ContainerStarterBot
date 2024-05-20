package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type ConfigStruct struct {
	Token                      string        `yaml:"token"`
	AppId                      string        `yaml:"appid"`
	Guild                      string        `yaml:"guild"`
	Channel                    string        `yaml:"channel"`
	Container                  string        `yaml:"container"`
	MaxTail                    int64         `yaml:"max_tail"`
	Timeout                    int           `yaml:"timeout"`
	EnableStop                 bool          `yaml:"enable_stop"`
	EnableAutoStop             bool          `yaml:"enable_auto_stop"`
	DefaultAutoStopDuration    time.Duration `yaml:"-"`
	DefaultAutoStopDurationRaw string        `yaml:"default_auto_stop_duration"`
	EnableAutoStopOverride     bool          `yaml:"enable_auto_stop_override"`
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
	if Config.EnableAutoStop {
		if Config.DefaultAutoStopDurationRaw != "off" {
			Config.DefaultAutoStopDuration, err = time.ParseDuration(Config.DefaultAutoStopDurationRaw)
			if err != nil {
				log.Fatal("Failed to parse duration of auto_stop_duration: ", err)
			}
		}
	}
}
