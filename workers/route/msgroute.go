package route

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Worker struct {
	MsgChan chan *tgbotapi.Update
}

func (w *Worker) Start() {
	for msg := range w.MsgChan {
		if msg.Message != nil {

		}
	}
}
