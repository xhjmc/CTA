package config

import (
	"encoding/json"
	"github.com/XH-JMC/cta/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	conf map[string]interface{}
}

func (c *Config) SetConf(conf map[string]interface{}) {
	c.conf = conf
}

func (c *Config) GetConf() map[string]interface{} {
	return c.conf
}

func (c *Config) LoadFromJSON(jsonBytes []byte) error {
	err := json.Unmarshal(jsonBytes, &c.conf)
	return err
}

func (c *Config) LoadFromYaml(yamlBytes []byte) error {
	if err := yaml.Unmarshal(yamlBytes, &c.conf); err != nil {
		return err
	}
	c.conf = util.ConvertMapInterfaceInterface2MapStringInterface(c.conf).(map[string]interface{})
	return nil
}

func (c *Config) LoadFromJSONFile(path string) error {
	jsonBytes, err := ioutil.ReadFile(path)
	if err == nil {
		err = c.LoadFromJSON(jsonBytes)
	}
	return err
}

func (c *Config) LoadFromYamlFile(path string) error {
	jsonBytes, err := ioutil.ReadFile(path)
	if err == nil {
		err = c.LoadFromYaml(jsonBytes)
	}
	return err
}

func (c *Config) Get(key string) interface{} {
	return c.conf[key]
}

func (c *Config) GetMap(key string) map[string]interface{} {
	ret, _ := c.conf[key].(map[string]interface{})
	return ret
}

func (c *Config) GetString(key string) string {
	str, _ := c.conf[key].(string)
	return str
}

func (c *Config) GetStringOrDefault(key string, defaultStr string) string {
	if str, ok := c.conf[key].(string); ok {
		return str
	}
	return defaultStr
}
