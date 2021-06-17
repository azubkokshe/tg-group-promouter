package config

import (
	"github.com/kelseyhightower/envconfig"
)

const (
	DataBase = "DB"
	BotAPi   = "TG"
)

// db configuration
type DBConfig struct {
	Host     string `required:"true"`
	Port     string `required:"true"`
	User     string `required:"true" default:"tgbot"`
	Password string `required:"true"`
	Name     string `required:"true" default:"tg_bot"`
}

// tg config
type TGConfig struct {
	Token string `required:"true"`
	Debug bool   `default:"true"`
}

type Config struct {
	DB DBConfig
	TG TGConfig
}

func ParseCfg() *Config {
	cfg := Config{}
	envconfig.MustProcess(DataBase, &cfg.DB)
	envconfig.MustProcess(BotAPi, &cfg.TG)
	return &cfg
}
