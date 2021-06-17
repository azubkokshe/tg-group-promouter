package route

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"sync"
)

type Route int

const (
	NewChannelRoute Route = 0
	NewUserInvite   Route = 1
)

func (r Route) String() string {
	switch r {
	case NewChannelRoute:
		return "Handle new channel"
	case NewUserInvite:
		return "Handle new members"
	}
	return "unknown"
}

type Worker struct {
	MsgChan chan *tgbotapi.Update
	Routes  map[Route]chan *tgbotapi.Update
	Wg      *sync.WaitGroup
}

func (w *Worker) Start() {
	defer w.Wg.Done()

	go func() {
		for msg := range w.MsgChan {

			fmt.Printf("receive: %#v\n", msg.Message)

			if isNewChannel(msg) {
				if err := w.send(NewChannelRoute, msg); err != nil {
					log.Println(err)
				}
				continue
			} else if isNewInvite(msg) {
				log.Println("WOW new invite")
				if err := w.send(NewUserInvite, msg); err != nil {
					log.Println(err)
				}
				continue
			}
		}
	}()
}

func (w *Worker) send(route Route, upd *tgbotapi.Update) error {
	if r, ok := w.Routes[route]; ok {
		r <- upd
		return nil
	}

	return fmt.Errorf("route not found for %s", route)
}

func isNewChannel(upd *tgbotapi.Update) bool {
	if upd.Message != nil {
		if upd.Message.ForwardFromChat != nil {
			return true
		}
	}

	return false
}

func isNewInvite(upd *tgbotapi.Update) bool {
	if upd.Message != nil {
		if upd.Message.NewChatMembers != nil && upd.Message.From != nil {
			return true
		}
	}

	return false
}
