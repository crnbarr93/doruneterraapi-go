package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Address  string `mapstructure:"address"`
	Database string `mapstructure:"database"`
	Testing  bool   `mapstructure:"testing"`
}

type Schema struct {
	Database DatabaseConfig `mapstructure:"database"`
	API      struct {
		Token string `mapstructure:"token"`
	} `mapstructure:"api"`
}

var (
	Config     *Schema
	TestConfig *Schema
)

func init() {
	config := viper.New()
	config.SetConfigName("config")
	config.AddConfigPath(".")
	config.AddConfigPath("config/")
	config.AddConfigPath("../config/")
	config.AddConfigPath("../")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	config.AutomaticEnv()

	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}
	err = config.Unmarshal(&Config)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}

	testConfig := viper.New()
	testConfig.SetConfigName("testconfig")
	testConfig.AddConfigPath(".")
	testConfig.AddConfigPath("config/")
	testConfig.AddConfigPath("../config/")
	// testConfig.AddConfigPath("../")
	testConfig.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	testConfig.AutomaticEnv()

	err = testConfig.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}
	err = testConfig.Unmarshal(&TestConfig)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}
}
