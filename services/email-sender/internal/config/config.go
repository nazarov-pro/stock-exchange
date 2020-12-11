package config

import (
	"github.com/spf13/viper"
	"strings"
	"os"
)

// Config - configuration context
var Config *viper.Viper = generate()

func generate() *viper.Viper {
	conf := viper.New()
	confFile := os.Getenv("CONFIG_FILE")
	if confFile == "" {
		confFile = `configs/app/config.yaml`
	}
	conf.SetConfigFile(confFile)
	err := conf.ReadInConfig()
	if err != nil {
		panic(err)
	}

	conf.SetEnvPrefix("_")
	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	conf.AutomaticEnv()
	return conf
}
