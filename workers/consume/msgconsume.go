package consume

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"sync"
)

type Worker struct {
	UpdatesChannel *tgbotapi.UpdatesChannel
	MsgChannel     chan tgbotapi.Update
	Wg             *sync.WaitGroup
}

func (w *Worker) Start() {
	defer w.Wg.Done()

	go func() {
		for update := range *w.UpdatesChannel {
			w.MsgChannel <- update
		}
	}()
}
