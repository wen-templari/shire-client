package core_test

import (
	"log"
	"testing"

	"github.com/templari/shire-client/core"
	"github.com/templari/shire-client/model"
)

var infoServerAddr = "http://localhost:3011"

func TestCoreLogin(t *testing.T) {
	core := core.MakeCore(infoServerAddr)
	_, err := core.Login(1, "12346")
	if err != nil {
		t.Error(err)
	}
}

func TestStartBobServer(t *testing.T) {
	bob := core.MakeCore(infoServerAddr)

	if _, err := bob.Login(7, "12345"); err != nil {
		t.Error(err)
	}

	for {
	}

}

func TestSendMessage(t *testing.T) {
	bob := core.MakeCore(infoServerAddr)

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

	sender := core.MakeCore(infoServerAddr)
	if _, err := sender.Register("tom", "12345"); err != nil {
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
func TestGetGroup(t *testing.T) {
	c := core.MakeCore(infoServerAddr)
	group, err := c.GetGroupById(1)
	if err != nil {
		log.Println(err)
	}
	log.Println(group)
}

func TestGroup(t *testing.T) {
	// a group of three users
	// alice, bob, tom
	alice := core.MakeCore(infoServerAddr)
	if _, err := alice.Register("alice", "12345"); err != nil {
		t.Error(err)
	}
	bob := core.MakeCore(infoServerAddr)
	if _, err := bob.Register("bob", "12345"); err != nil {
		t.Error(err)
	}
	tom := core.MakeCore(infoServerAddr)
	if _, err := tom.Register("tom", "12345"); err != nil {
		t.Error(err)
	}

	// alice creates a group
	group, err := alice.CreateGroup([]int{alice.GetUser().Id, bob.GetUser().Id, tom.GetUser().Id})
	if err != nil {
		t.Error(err)
	}

	alice.SendMessage(model.Message{
		GroupId: group.Id,
		From:    alice.GetUser().Id,
		Content: "hello From Alice",
		Time:    "123",
	})

	alice.SendMessage(model.Message{
		GroupId: group.Id,
		From:    alice.GetUser().Id,
		Content: "hello Again From Alice",
		Time:    "123",
	})
	for {
	}

}
