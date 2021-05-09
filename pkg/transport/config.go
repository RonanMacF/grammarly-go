package transport

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Dialects[]string`mapstructure:"dialect"`
}

func (c Config) VerifyConfig() error {
	if err := verifyDialect(c.Dialects); err != nil{
		return err
	}
	return nil
}

func verifyDialect(dialects []string)error{
	allowedDialects := map[string]struct{}{
		"american" : {},
		"british" :{},
		"canadian": {},
		"australian":{},
	}
	for _, d := range dialects{
		if _, ok := allowedDialects[strings.ToLower(d)]; !ok{
			return fmt.Errorf("invalid dialect %s, allowed dialects are american, british, " +
				"canadian and australian", d)
		}
	}
	return nil
}

func (c Config) GenerateConfigMessage() map[string]interface{} {
	fmt.Println(strings.Join(c.Dialects, " | "))
	return map[string]interface{}{
		"type":            "initial",
		"docid":           "1234",
		"client":          "extension_chrome",
		"protocolVersion": "1.0",
		"clientSupports": []string{
			"free_clarity_alerts",
			"readability_check",
			"filler_words_check",
			"free_occasional_premium_alerts",
		},
		"dialect":       strings.Join(c.Dialects, " | "),
		//"dialect":       "american",
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
