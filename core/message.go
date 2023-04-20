package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/templari/shire-client/core/model"
)

func (c *Core) ReceiveMessage(message model.Message) error {
	//TODO not implemented
	// save message to db ??
	// notify watchers
	for _, ch := range c.subscribers {
		ch <- message
	}
	return nil
}

func (c *Core) SendMessage(message model.Message) error {
	//TODO not implemented
	// find message.to's address
	receiver, err := c.GetUserById(message.To)
	if err != nil {
		return err
	}
	if receiver.Address == "" || receiver.Port == 0 {
		return fmt.Errorf("receiver's address is not set")
	}
	// send message to address
	bytesData, err := json.Marshal(message)
	if err != nil {
		return err
	}
	requestURL := fmt.Sprintf("http://%v:%v/message", receiver.Address, receiver.Port)
	req, err := http.NewRequest("POST", requestURL, bytes.NewReader(bytesData))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message to %v:%v", receiver.Address, receiver.Port)
	}
	// save message to db ???
	// notify watchers
	for _, ch := range c.subscribers {
		ch <- message
	}
	return nil
}

func (c *Core) Subscribe(ch chan model.Message) {
	c.subscribers = append(c.subscribers, ch)
}
