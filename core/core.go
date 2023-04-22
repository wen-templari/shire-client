package core

import (
	"github.com/templari/shire-client/model"
)

type Core struct {
	InfoServerAddress string
	user              model.User
	token             string
	subscribers       []chan model.Message
}

func MakeCore(address string) *Core {
	return &Core{
		InfoServerAddress: address,
		subscribers:       make([]chan model.Message, 0),
	}
}

func (c *Core) Ping() string {
	return "pong"
}
