package core

import "github.com/templari/shire-client/core/model"

func (c *Core) ReceiveMessage(message model.Message) error {
	//TODO not implemented
	// save message to db
	// notify watchers
	return nil
}

func (c *Core) SendMessage(message model.Message) error {
	//TODO not implemented
	// find message.to's address
	// send message to address
	// save message to db
	return nil
}
