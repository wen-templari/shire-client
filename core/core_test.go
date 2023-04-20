package core_test

import (
	"log"
	"testing"

	"github.com/templari/shire-client/core"
	"github.com/templari/shire-client/core/model"
)

func TestCoreLogin(t *testing.T) {
	core := core.MakeCore("http://localhost:3011", nil)
	_, err := core.Login(1, "12346")
	if err != nil {
		t.Error(err)
	}
}

func TestSendMessage(t *testing.T) {
	bob := core.MakeCore("http://localhost:3011", nil)

	if _, err := bob.Register("bob", "12345"); err != nil {
		t.Error(err)
	}

	messageChan := make(chan model.Message)
	bob.Subscribe(messageChan)

	go func() {
		for {
			message := <-messageChan
			log.Printf("received message: %v", message)
		}
	}()

	sender := core.MakeCore("http://localhost:3011", nil)
	if _, err := bob.Register("tom", "12345"); err != nil {
		t.Error(err)
	}

	receiver, err := sender.GetUserById(bob.GetUser().Id)
	if err != nil {
		t.Error(err)
	}

	err = sender.SendMessage(model.Message{
		To:      receiver.Id,
		From:    sender.GetUser().Id,
		Content: "hello",
		Time:    "123",
	})
	if err != nil {
		t.Errorf("failed to send message: %v", err)
	}

}
