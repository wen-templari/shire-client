package main

import (
	"context"
	"fmt"

	"github.com/templari/shire-client/core"
	"github.com/templari/shire-client/model"
)

// App struct
type App struct {
	ctx context.Context

	messageChan chan model.Message
	core        *core.Core
}

// NewApp creates a new App application struct
func NewApp() *App {

	app := &App{
		core:        core.MakeCore("http://localhost:3011"),
		messageChan: make(chan model.Message),
	}

	app.core.Subscribe(app.messageChan)

	return app

}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// if err := godotenv.Load(); err != nil {
	// 	util.Logger.Fatal("Error loading .env file")
	// }

	go a.messageUpdateHandler()

}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
