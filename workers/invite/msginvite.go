package invite

import (
	"context"
	"github.com/azubkokshe/tg-group-promouter/models"
	"github.com/azubkokshe/tg-group-promouter/store/channels"
	"github.com/azubkokshe/tg-group-promouter/store/invites"
	"github.com/azubkokshe/tg-group-promouter/store/users"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"log"
	"sync"
)

type Worker struct {
	MsgChannel   chan *tgbotapi.Update
	Wg           *sync.WaitGroup
	Bot          *tgbotapi.BotAPI
	ChannelStore channels.Store
	UserStore    users.Store
	InviteStore  invites.Store
	DB           *sqlx.DB
}

func (w *Worker) Start() {
	defer w.Wg.Done()

	go func() {
		for update := range w.MsgChannel {
			//Сначала узнаем мониторим ли мы канал
			chnl, err := w.ChannelStore.GetByID(context.Background(), update.Message.Chat.ID)
			if err != nil {
				log.Println("channel to handle new member not found:", update.Message.Chat.ID, err)
				continue
			}

			u2a := make([]tgbotapi.User, 0, len(*update.Message.NewChatMembers)+1)
			u2a = append(u2a, *update.Message.From)
			u2a = append(u2a, *update.Message.NewChatMembers...)

			tx := w.DB.MustBegin()

			//Теперь зарегаем нового юзера, который пригласил
			for _, u := range u2a {
				if err = w.UserStore.Store(context.Background(), tx, &models.User{
					ID:        int64(u.ID),
					FirstName: u.FirstName,
					LastName:  u.LastName,
					IsBot:     u.IsBot,
					UserName:  u.UserName,
				}); err != nil {
					log.Println(err)
					break
				}
			}

			if err != nil {
				if err := tx.Rollback(); err != nil {
					log.Println(err)
				}
				continue
			}

			for _, m := range *(update.Message.NewChatMembers) {
				if err = w.InviteStore.Store(context.Background(), tx, &models.Invites{
					ChannelID: chnl.ID,
					FromID:    int64(update.Message.From.ID),
					MemberID:  int64(m.ID),
				}); err != nil {
					log.Println(err)
					break
				}
			}

			if err != nil {
				if err := tx.Rollback(); err != nil {
					log.Println(err)
				}
				continue
			}

			if err := tx.Commit(); err != nil {
				log.Println(err)
			}
		}
	}()
}
