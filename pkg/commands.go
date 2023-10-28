package pkg

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Commands struct {
	Commands []Command `yaml:"Commands"`
}

type Command struct {
	Name string `yaml:"name"`
	URI  string `yaml:"uri"`
}

func (c *Commands) ReadCommands() error {
	yamlFile, err := os.ReadFile("./conf/commands.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return err
	}
	return nil
}

func (cs *Commands) GetCommand(s string) (*Command, error) {
	for _, c := range cs.Commands {
		if c.Name == s {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("could not find command: %s, please check the 'commands' file", s)
}
