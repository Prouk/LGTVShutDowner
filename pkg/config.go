package pkg

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	LGTVShutDowner struct {
		TVInfos struct {
			Ip  string `yaml:"Ip"`
			Mac string `yaml:"Mac"`
		} `yaml:"TVInfos"`
	} `yaml:"LGTVShutDowner"`
}

func (c *Config) ReadConf() error {
	yamlFile, err := os.ReadFile("./conf/config.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return err
	}
	return nil
}
