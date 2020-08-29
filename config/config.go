package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

// PlatformSettings spec
type PlatformSettings struct {
	Config struct {
		KeyCloak struct {
			Realm     string `yaml:"realm"`
			Namespace string `yaml:"namespace"`
			FQDN      string `yaml:"fqdn"`
		}
	}
}

// Configuration spec
type Configuration struct {
	Tenant TenantConfig
}

// TenantConfig spec
type TenantConfig struct {
	Name   string `yaml:"name"`
	Groups []GroupConfig
}

// GroupConfig spec
type GroupConfig struct {
	Name    string   `yaml:"name"`
	Admin   bool     `yaml:"admin"`
	Members []string `yaml:"members"`
}

// LoadTenantConfig reads config from YAML
func LoadTenantConfig(fileName string) TenantConfig {

	var config Configuration
	readConfig(fileName, &config)

	return config.Tenant
}

// LoadPlatformSettings reads platform settings from YAML
func LoadPlatformSettings(fileName string) PlatformSettings {

	var config PlatformSettings
	readConfig(fileName, &config)

	return config
}

func readConfig(fileName string, config interface{}) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error reading YAML file: %s\n", err)
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Fatalf("Error parsing YAML file: %s\n", err)
	}
}
