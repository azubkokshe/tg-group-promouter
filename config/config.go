package config

import (
	"github.com/kelseyhightower/envconfig"
)

const (
	DataBase = "DB"
	BotAPi   = "TG"
)

// db configuration
type dbConfig struct {
	Host     string `required:"true"`
	Port     string `required:"true"`
	User     string `required:"true" default:"tgbot"`
	Password string `required:"true"`
	Name     string `required:"true" default:"tg_bot"`
}

// tg config
type tgConfig struct {
	Token string `required:"true"`
	Debug bool   `default:"true"`
}

type config struct {
	DB dbConfig
	TG tgConfig
}

func ParseCfg() *config {
	cfg := config{}
	envconfig.MustProcess(DataBase, &cfg.DB)
	envconfig.MustProcess(BotAPi, &cfg.TG)
	return &cfg
}
