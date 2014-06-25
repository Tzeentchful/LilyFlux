package config

import "io/ioutil"
import yaml "launchpad.net/goyaml"

type Config struct {
	Connect ConfigConnect `yaml:"connect"`
	Influx ConfigInflux `yaml:"Influx"`
}

type ConfigInflux struct {
	Address string `yaml:address`
	Port uint16 `yaml:port`
	Username string `yaml:username`
	Password string `yaml:password`
	Database string `yaml:database`
}

type ConfigConnect struct {
	Address string `yaml:"address"`
	Credentials ConfigConnectCredentials `yaml:"credentials"`
}

type ConfigConnectCredentials struct {
	Username string
	Password string
}


func DefaultConfig() (config *Config) {
	return &Config{
		Connect: ConfigConnect {
			Address: "127.0.0.1:5091",
			Credentials: ConfigConnectCredentials{
				Username: "example",
				Password: "example",
			},
		},
		Influx: ConfigInflux {
			Address: "localhost",
			Port: 8086,
			Username: "root",
			Password: "root",
			Database: "lilyflux",
		},
	}
}

func LoadConfig(file string) (config *Config, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	var cfg Config
	config = &cfg
	err = yaml.Unmarshal(data, config)
	return
}

func SaveConfig(file string, config *Config) (err error) {
	data, err := yaml.Marshal(config)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(file, data, 0644)
	return
}
