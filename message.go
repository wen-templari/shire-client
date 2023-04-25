package main

import (
	"log"
	"time"

	"github.com/templari/shire-client/model"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a App) Ping() string {
	return a.core.Ping()
}

func (a *App) SendMessage(message model.Message) error {
	message.Time = time.Now().Format(time.RFC3339)
	log.Printf("Sending message: %v", message)
	return a.core.SendMessage(message)
}

// for testing only, this should not be needed
func (a *App) ReceiveMessage(message model.Message) error {
	log.Printf("Received message: %v", message)
	return a.core.ReceiveMessage(message)
}

// start this as a go routine at startup
func (a *App) messageUpdateHandler() {
	for {
		message := <-a.messageChan
		runtime.EventsEmit(a.ctx, "onMessage", message)

		// TODO  persist
	}
}
