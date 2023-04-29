package main

import (
	"log"

	"github.com/templari/shire-client/model"
	"github.com/templari/shire-client/repo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a App) Ping() string {
	return a.core.Ping()
}

func (a *App) SendMessage(message model.Message) error {
	err := a.core.SendMessage(message)
	if err != nil {
		log.Printf(" err Sending message: %v, err: %v", message, err)
	}
	return err
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

		err := repo.SaveMessage(&message)
		if err != nil {
			log.Println("error on persisting data", err)
		}
	}
}

func (a *App) GetMessages() ([]model.Message, error) {
	return repo.GetAllMessage()
}
