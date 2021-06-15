package main

import (
	"github.com/azubkokshe/tg-group-promouter/config"
	"github.com/azubkokshe/tg-group-promouter/workers/consume"
	route "github.com/azubkokshe/tg-group-promouter/workers/route"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
)

func init() {
	os.Setenv("DB_HOST", "127.0.0.1")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "tgbot")
	os.Setenv("DB_PASSWORD", "tg_b_o_t")
	os.Setenv("DB_NAME", "tg_bot")

	os.Setenv("TG_TOKEN", "1855914325:AAHqPWswWYbgLA-v3ue8sA0e6au_4wvQBZI")
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	cfg := config.ParseCfg()

	m, err := migrate.New("file://migrations", "postgres://"+cfg.DB.User+":"+
		cfg.DB.Password+"@"+cfg.DB.Host+":"+cfg.DB.Port+"/"+cfg.DB.Name+"?sslmode=disable")
	panicOnErr(err)

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TG.Token)
	panicOnErr(err)
	bot.Debug = cfg.TG.Debug
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	panicOnErr(err)

	mc := make(chan *tgbotapi.Update)

	consumer := consume.Worker{UpdatesChannel: &updates, MsgChannel: mc}
	consumer.Start()

	router := route.Worker{MsgChan: mc}
	router.Start()

	log.Printf("Authorized on account %s", bot.Self.UserName)
}
