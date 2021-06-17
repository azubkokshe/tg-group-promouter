package channel

import (
	"context"
	"github.com/azubkokshe/tg-group-promouter/models"
	"github.com/azubkokshe/tg-group-promouter/store/channels"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"log"
	"sync"
)

type Worker struct {
	MsgChannel chan *tgbotapi.Update
	Wg         *sync.WaitGroup
	Bot        *tgbotapi.BotAPI
	Store      channels.Store
	DB         *sqlx.DB
}

func (w *Worker) Start() {
	defer w.Wg.Done()

	go func() {
		for update := range w.MsgChannel {
			// получаем список админов канала, с которого нам прислали сообщение
			chMembers, err := w.Bot.GetChatAdministrators(update.Message.ForwardFromChat.ChatConfig())
			if err != nil {
				log.Println("an error occurred", err)
				continue
			}

			cond := 2

			for _, m := range chMembers {
				// sender must be admin of channel
				if m.User.ID == update.Message.From.ID {
					cond--
					// current bot must be admin of channel
				} else if m.User.IsBot && m.User.ID == w.Bot.Self.ID {
					cond--
				}
			}

			if cond != 0 {
				log.Println("channel can't be added to channel list. One of the main conditions is not met")
				continue
			}

			tx := w.DB.MustBegin()

			if err := w.Store.Store(context.Background(), tx, &models.Channel{
				ID:    update.Message.ForwardFromChat.ID,
				Title: update.Message.ForwardFromChat.Title,
			}); err != nil {
				log.Println("channel can't be added to channel list", err)
				if err := tx.Rollback(); err != nil {
					log.Println(err)
				}
			} else {
				if err := tx.Commit(); err != nil {
					log.Println(err)
				}
			}
		}
	}()
}
