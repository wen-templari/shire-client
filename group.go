package main

import "github.com/templari/shire-client/model"

func (a *App) CreateGroup(idList []int) (model.Group, error) {
	return a.core.CreateGroup(idList)
}

func (a *App) GetGroupById(id int) (model.Group, error) {
	return a.core.GetGroupById(id)
}
