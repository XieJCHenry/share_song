package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MySql *MySqlConfig `json:"mysql,omitempty" yaml:"mysql"`
	Mongo *MongoConfig `json:"mongo,omitempty" yaml:"mongo"`
}

type MySqlConfig struct {
	Host     string `json:"host,omitempty" yaml:"host"`
	Port     int    `json:"port,omitempty" yaml:"port"`
	UserName string `json:"userName,omitempty" yaml:"username"`
	Password string `json:"password,omitempty" yaml:"password"`
	DataBase string `json:"dataBase" yaml:"database"`
	Timeout  string `json:"timeout" yaml:"timeout"`
}

type MongoConfig struct {
	Host string `json:"host,omitempty" yaml:"host"`
	Port int    `json:"port,omitempty" yaml:"port"`
}

func loadConfig() (*Config, error) {

	curPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}

	// D:\workspace\go\share_song\cmd\conf\conf.yaml
	configPath := filepath.Join(curPath, "..", "..", "conf", "conf.yaml")
	all, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(all, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
