package consume

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Worker struct {
	UpdatesChannel *tgbotapi.UpdatesChannel
	MsgChannel     chan *tgbotapi.Update
}

func (w *Worker) Start() {
	go func() {
		for update := range *w.UpdatesChannel {

			if update.Message.NewChatMembers != nil &&
				len(*(update.Message.NewChatMembers)) > 0 &&
				update.Message.From != nil {
			}
		}
	}()
}
