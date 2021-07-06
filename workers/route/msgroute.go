package route

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"sync"
)

type Route int

const (
	UnknownRoute    Route = -1
	NewChannelRoute Route = 0
	NewUserInvite   Route = 1
	NewCommand      Route = 2
)

func (r Route) String() string {
	switch r {
	case UnknownRoute:
		return "Unknown route"
	case NewChannelRoute:
		return "Handle new channel registration on bot"
	case NewUserInvite:
		return "Handle new members in channel"
	case NewCommand:
		return "Handle new bot command"
	}
	return "unknown"
}

type Worker struct {
	MsgChan chan tgbotapi.Update
	Routes  map[Route]chan tgbotapi.Update
	Wg      *sync.WaitGroup
}

func (w *Worker) Start() {
	defer w.Wg.Done()

	go func() {
		for msg := range w.MsgChan {
			r := UnknownRoute

			if isNewChannel(msg) {
				r = NewChannelRoute
			} else if isNewInvite(msg) {
				r = NewUserInvite
			} else if msg.Message.IsCommand() {
				r = NewCommand
			}

			if r != UnknownRoute {
				if err := w.send(r, msg); err != nil {
					log.Println(err)
				}
			}
		}
	}()
}

func (w *Worker) send(route Route, upd tgbotapi.Update) error {
	if r, ok := w.Routes[route]; ok {
		fmt.Println("send to route", route.String())
		r <- upd
		return nil
	}

	return fmt.Errorf("route not found for %s", route)
}

func isNewChannel(upd tgbotapi.Update) bool {
	if upd.Message != nil {
		if upd.Message.ForwardFromChat != nil {
			return true
		}
	}

	return false
}

func isNewInvite(upd tgbotapi.Update) bool {
	if upd.Message != nil {
		if upd.Message.NewChatMembers != nil && upd.Message.From != nil {
			return true
		}
	}

	return false
}
