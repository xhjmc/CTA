package config

var config = &Config{}

func SetConf(conf map[string]interface{}) {
	config.SetConf(conf)
}

func GetConf() map[string]interface{} {
	return config.GetConf()
}

func LoadFromJSON(jsonBytes []byte) error {
	return config.LoadFromJSON(jsonBytes)
}

func LoadFromYaml(yamlBytes []byte) error {
	return config.LoadFromYaml(yamlBytes)
}

func LoadFromJSONFile(path string) error {
	return config.LoadFromJSONFile(path)
}

func LoadFromYamlFile(path string) error {
	return config.LoadFromYamlFile(path)
}

func Get(key string) interface{} {
	return config.Get(key)
}

func GetMap(key string) map[string]interface{} {
	return config.GetMap(key)
}

func GetString(key string) string {
	return config.GetString(key)
}

func GetStringOrDefault(key string, defaultStr string) string {
	return config.GetStringOrDefault(key, defaultStr)
}
