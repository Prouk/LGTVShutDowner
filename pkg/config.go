package pkg

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	LGTVShutDowner struct {
		TVInfos struct {
			Ip        string `yaml:"ip"`
			Mac       string `yaml:"mac"`
			ClientKey string `yaml:"clientKey"`
		} `yaml:"TVInfos"`
	} `yaml:"LGTVShutDowner"`
}

func (lsd *Lsd) LoadConfig() {
	file, err := os.ReadFile(lsd.ConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(file, lsd.Config)
	if err != nil {
		log.Fatal(err)
	}
}

func (lsd *Lsd) SaveConfig() {
	data, err := yaml.Marshal(lsd.Config)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(lsd.ConfigFilePath, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
