package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/templari/shire-client/model"
)

func (c *Core) ReceiveMessage(message model.Message) error {
	// notify watchers
	log.Printf("%v: Received message: %v", c.user.Id, message)
	for _, ch := range c.subscribers {
		ch <- message
		log.Println(len(c.subscribers))
	}
	return nil
}

func (c *Core) SendMessage(message model.Message) error {
	message.Time = time.Now().Format(time.RFC3339)
	if message.GroupId <= 0 {
		log.Println("one to one")
		return c.sendOneToOneMessage(message)
	} else {
		log.Println("group message")
		return c.sendGroupMessage(message)
	}
}

func (c *Core) sendOneToOneMessage(message model.Message) error {
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
	// notify watchers
	for _, ch := range c.subscribers {
		ch <- message
	}
	log.Printf("%v: Sending message: %v", c.user.Id, message)
	return nil
}

func (c *Core) sendGroupMessage(message model.Message) error {
	// find group's wrapper
	w, ok := c.wrappers[message.GroupId]

	if !ok {
		_, err := c.SetupGroup(message.GroupId)
		if err != nil {
			log.Printf("wrapper not found for group %v", message.GroupId)
			return fmt.Errorf("wrapper not found for group %v", message.GroupId)
		}
	}

	s, err := json.Marshal(message)
	if err != nil {
		return err
	}
	w.client.Append(strconv.Itoa(message.GroupId), string(s))

	log.Printf("%v: Sending message: %v", c.user.Id, message)
	return nil
}

func (c *Core) Subscribe(ch chan model.Message) {
	c.subscribers = append(c.subscribers, ch)
}
