package main

import (
	"fmt"
	"github.com/azubkokshe/tg-group-promouter/config"
	"github.com/azubkokshe/tg-group-promouter/store/channels/channels_pg"
	"github.com/azubkokshe/tg-group-promouter/store/invites/invites_pg"
	"github.com/azubkokshe/tg-group-promouter/store/users/users_pg"
	"github.com/azubkokshe/tg-group-promouter/workers/channel"
	"github.com/azubkokshe/tg-group-promouter/workers/command"
	"github.com/azubkokshe/tg-group-promouter/workers/consume"
	"github.com/azubkokshe/tg-group-promouter/workers/invite"
	"github.com/azubkokshe/tg-group-promouter/workers/route"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/lib/pq"
)

func init() {
	//os.Setenv("DB_HOST", "127.0.0.1")
	//os.Setenv("DB_PORT", "5432")
	//os.Setenv("DB_USER", "tgbot")
	//os.Setenv("DB_PASSWORD", "tg_b_o_t")
	//os.Setenv("DB_NAME", "tg_bot")
	//
	//os.Setenv("TG_TOKEN", "1855914325:AAHqPWswWYbgLA-v3ue8sA0e6au_4wvQBZI")
	//os.Setenv("TG_DEBUG", "true")
}

//TODO update channel title if exists
func main() {
	cfg := config.ParseCfg()

	err := doMigrations(cfg.DB)
	panicOnErr(err)

	db, err := doDBConnect(cfg.DB)
	panicOnErr(err)

	bot, err := tgbotapi.NewBotAPI(cfg.TG.Token)
	panicOnErr(err)
	bot.Debug = cfg.TG.Debug
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	panicOnErr(err)

	routes := make(map[route.Route]chan tgbotapi.Update)
	routes[route.NewChannelRoute] = make(chan tgbotapi.Update, 10000)
	routes[route.NewUserInvite] = make(chan tgbotapi.Update, 10000)
	routes[route.NewCommand] = make(chan tgbotapi.Update, 10000)

	mc := make(chan tgbotapi.Update)
	termChan := make(chan os.Signal, 1)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	consumer := consume.Worker{
		UpdatesChannel: &updates,
		MsgChannel:     mc,
		Wg:             wg,
	}
	consumer.Start()

	wg.Add(1)
	router := route.Worker{
		MsgChan: mc,
		Wg:      wg,
		Routes:  routes,
	}
	router.Start()

	wg.Add(1)
	channels := channel.Worker{
		MsgChannel: routes[route.NewChannelRoute],
		Wg:         wg,
		Bot:        bot,
		DB:         db,
		Store:      channels_pg.NewRepository(db),
	}
	channels.Start()

	wg.Add(1)
	invites := invite.Worker{
		MsgChannel:   routes[route.NewUserInvite],
		Wg:           wg,
		Bot:          bot,
		DB:           db,
		UserStore:    users_pg.NewRepository(db),
		InviteStore:  invites_pg.NewRepository(db),
		ChannelStore: channels_pg.NewRepository(db),
	}
	invites.Start()

	wg.Add(1)
	commands := command.Worker{
		MsgChannel:  routes[route.NewCommand],
		Wg:          wg,
		Bot:         bot,
		DB:          db,
		InviteStore: invites_pg.NewRepository(db),
	}
	commands.Start()

	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)

	log.Printf("Authorized on account %s", bot.Self.UserName)

	<-termChan

	bot.StopReceivingUpdates()
	close(mc)

	for _, v := range routes {
		close(v)
	}

	wg.Wait()
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func doMigrations(cfg config.DBConfig) error {
	m, err := migrate.New("file://migrations", "postgres://"+cfg.User+":"+
		cfg.Password+"@"+cfg.Host+":"+cfg.Port+"/"+cfg.Name+"?sslmode=disable")

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func doDBConnect(cfg config.DBConfig) (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}
