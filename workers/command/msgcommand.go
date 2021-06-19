package command

import (
	"context"
	"fmt"
	"github.com/azubkokshe/tg-group-promouter/store/invites"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	"strings"
	"sync"
)

type Worker struct {
	MsgChannel  chan *tgbotapi.Update
	Wg          *sync.WaitGroup
	Bot         *tgbotapi.BotAPI
	InviteStore invites.Store
	DB          *sqlx.DB
}

func (w *Worker) Start() {
	defer w.Wg.Done()

	go func() {
		for update := range w.MsgChannel {
			switch update.Message.Command() {
			case "rating":
				rating := w.InviteStore.SelectRating(context.Background(), update.Message.Chat.ID)

				text := strings.Builder{}

				text.WriteString("Рейтинг участников по количеству пригласивших (10 позиций): \n\n")

				for i, r := range *rating {
					if i > 9 {
						break
					}
					if i > 0 {
						text.WriteString("\n")
					}
					text.WriteString(fmt.Sprintf("%d", i+1))
					text.WriteString(". ")
					text.WriteString(fmt.Sprintf("[%s](tg://user?id=%d) (%d)", r.UserName, r.UserID, r.Count))
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text.String())
				msg.ReplyToMessageID = update.Message.MessageID
				msg.ParseMode = "Markdown"

				_, err := w.Bot.Send(msg)
				if err != nil {
					fmt.Println("an error occurred", err)
				}
			case "position":
				rating := w.InviteStore.SelectRating(context.Background(), update.Message.Chat.ID)

				idx := -1

				for i, r := range *rating {
					if r.UserID == int64(update.Message.From.ID) {
						idx = i
						break
					}
				}

				if idx < 0 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Для Вас отсутствует запись в турнирной таблице. Пригласите кого-нибудь =)")
					msg.ReplyToMessageID = update.Message.MessageID
					_, err := w.Bot.Send(msg)
					if err != nil {
						fmt.Println("an error occurred", err)
					}
					continue
				}

				text := strings.Builder{}

				text.WriteString(fmt.Sprintf("Количество приглашенных: %d\n", (*rating)[idx].Count))
				text.WriteString(fmt.Sprintf("Место в рейтинге: %d из %d", idx+1, len(*rating)))

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text.String())
				msg.ReplyToMessageID = update.Message.MessageID
				_, err := w.Bot.Send(msg)
				if err != nil {
					fmt.Println("an error occurred", err)
				}
			default:
				fmt.Println("Unknown bot command", update.Message.Command())
			}
		}
	}()
}
