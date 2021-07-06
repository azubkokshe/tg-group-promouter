package invite

import (
	"context"
	"encoding/json"
	"fmt"
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
	MsgChannel   chan tgbotapi.Update
	Wg           *sync.WaitGroup
	Bot          *tgbotapi.BotAPI
	ChannelStore channels.Store
	UserStore    users.Store
	InviteStore  invites.Store
	DB           *sqlx.DB
}

func (w *Worker) regInvite(update tgbotapi.Update) error {
	//Сначала узнаем мониторим ли мы канал
	chnl, err := w.ChannelStore.GetByID(context.Background(), update.Message.Chat.ID)
	if err != nil {
		return err
	}

	u2a := make([]tgbotapi.User, 0, len(*update.Message.NewChatMembers)+1)
	u2a = append(u2a, *update.Message.From)
	u2a = append(u2a, *update.Message.NewChatMembers...)

	tx := w.DB.MustBegin()

	//Необходимо на всякие пожарные на первых порах сохранять записи в журнал
	if j, err := json.Marshal(update); err == nil {
		if err := w.InviteStore.Journal(context.Background(), tx, &models.Journal{
			Record: j,
		}); err != nil {
			if err := tx.Rollback(); err != nil {
				log.Println(err)
			}
			return err
		}
	} else {
		if err := tx.Rollback(); err != nil {
			log.Println(err)
		}
		return err
	}

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
		return err
	}

	for _, m := range *(update.Message.NewChatMembers) {
		if m.IsBot {
			continue
		}

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
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (w *Worker) Start() {
	defer w.Wg.Done()

	go func() {
		for update := range w.MsgChannel {
			fmt.Println("new invite")
			if err := w.regInvite(update); err != nil {
				log.Println(err)
			}
		}
	}()
}
