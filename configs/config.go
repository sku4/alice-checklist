package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port       string `mapstructure:"port"`
	DBPath     string `mapstructure:"db_file"`
	GoogleKeep `mapstructure:"googlekeep"`
}

type GoogleKeep struct {
	NoteRootId string `mapstructure:"noteRootId"`
}

func Init() (*Config, error) {
	mainViper := viper.New()
	mainViper.AddConfigPath("configs")
	if err := mainViper.ReadInConfig(); err != nil {
		return nil, err
	}

	gkViper := viper.New()
	gkViper.AddConfigPath("configs/googlekeep")
	if err := gkViper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config

	if err := mainViper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	if err := gkViper.Unmarshal(&cfg.GoogleKeep); err != nil {
		return nil, err
	}

	return &cfg, nil
}
