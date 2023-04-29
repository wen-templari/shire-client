package main

import (
	"log"

	"github.com/templari/shire-client/model"
	"github.com/templari/shire-client/repo"
)

func (a *App) Register(name string, password string) (model.User, error) {
	u, err := a.core.Register(name, password)
	if err != nil {
		log.Println(err)
	}
	repo.InitDB(u.Id)
	return u, err

}

func (a *App) Login(id int, password string) (model.User, error) {
	u, err := a.core.Login(id, password)
	if err != nil {
		log.Println(err)
	}
	repo.InitDB(u.Id)
	return u, err
}

func (a *App) UpdateUser() (model.User, error) {
	return a.core.UpdateUser()
}

func (a *App) GetUsers() ([]model.User, error) {
	return a.core.GetUsers()
}

func (a *App) GetUserById(id int) (model.User, error) {
	return a.core.GetUserById(id)
}
