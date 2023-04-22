package main

import (
	"github.com/templari/shire-client/model"
)

func (a *App) Register(name string, password string) (model.User, error) {
	return a.core.Register(name, password)
}

func (a *App) Login(id int, password string) (model.User, error) {
	return a.core.Login(id, password)
}

func (a *App) UpdateUser(port int) (model.User, error) {
	return a.core.UpdateUser(port)
}

func (a *App) GetUsers() ([]model.User, error) {
	return a.core.GetUsers()
}

func (a *App) GetUserById(id int) (model.User, error) {
	return a.core.GetUserById(id)
}
