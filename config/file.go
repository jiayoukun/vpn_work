package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type file struct {
	path      string
	cfile      string
}

func (f *file)load() (*Config, error) {
	configFile := f.path
	f.cfile = configFile
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	c := &Config{}

	err = yaml.Unmarshal(content, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}


