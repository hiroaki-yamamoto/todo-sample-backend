package main

import "github.com/spf13/viper"

type Config struct {
	Port   string
	DB     string
	Secret string
}

func LoadConfig(cfgName string) (Config, error) {
	var cfg Config
	viper.SetConfigName(cfgName)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
