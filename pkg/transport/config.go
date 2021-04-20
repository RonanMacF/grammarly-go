package transport

import (
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Dialect[]string`mapstructure:"dialect"`
}

func (c Config) GenerateConfigMessage() map[string]interface{} {
	return map[string]interface{}{
		"type":            "initial",
		"docid":           "1234",
		"client":          "extension_chrome",
		"protocolVersion": "1.0",
		"clientSupports": []string{
			"free_clarity_alerts",
			"readability_check",
			"filler_words_check",
			"sentence_variety_check",
			"free_occasional_premium_alerts",
		},
		"dialect":       strings.Join(c.Dialect, " | "),
		"clientVersion": "14.924.2437",
		"extDomain":     "keep.google.com",
		"action":        "start",
		"id":            "0",
		"sid":           "0",
	}
}

func LoadConfig(path string) (Config, error) {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			//fmt.Infof("no configuration file called %s.%s found in %s", configFileName,
			//	configFiletype, path)
			return Config{}, nil
		}
		return Config{}, err
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil{
		return Config{}, err
	}
	return c, nil
}
