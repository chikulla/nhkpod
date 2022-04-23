package nhkpod

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Conf struct {
	Podcasts []struct {
		ID       string `yaml:"id"`
		CornerID string `yaml:"corner_id"`
	} `yaml:"podcasts"`
}

func GetConf(confPath string) (*Conf, error) {
	file, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	var c Conf
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
