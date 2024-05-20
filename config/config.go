package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type ConfigStruct struct {
	Token                           string        `yaml:"token"`
	AppId                           string        `yaml:"appid"`
	Guild                           string        `yaml:"guild"`
	Channel                         string        `yaml:"channel"`
	Container                       string        `yaml:"container"`
	MaxTail                         int64         `yaml:"max_tail"`
	Timeout                         int           `yaml:"timeout"`
	EnableStop                      bool          `yaml:"enable_stop"`
	EnableAutoShutdown              bool          `yaml:"enable_auto_shutdown"`
	DefaultAutoShutdownDuration     time.Duration `yaml:"-"`
	DefaultAutoShutdownDuration_raw string        `yaml:"default_auto_shutdown_duration"`
	EnableAutoShutdownOverride      bool          `yaml:"enable_auto_shutdown_override"`
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
	if Config.DefaultAutoShutdownDuration_raw != "off" {
		Config.DefaultAutoShutdownDuration, err = time.ParseDuration(Config.DefaultAutoShutdownDuration_raw)
		if err != nil {
			log.Fatal("Failed to parse duration of auto_shutdown_duration: ", err)
		}
	}

}
