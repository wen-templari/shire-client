package repo_test

import (
	"log"
	"testing"
	"time"

	"github.com/templari/shire-client/model"
	"github.com/templari/shire-client/repo"
)

func TestMain(m *testing.M) {
	repo.InitDB(1)
	m.Run()
}

func TestSaveMessage(t *testing.T) {
	message := model.Message{
		From:    1,
		To:      1,
		Content: "hello",
		GroupId: 0,
		Time:    time.Now().Format(time.RFC3339),
	}
	log.Println(message)
	if err := repo.SaveMessage(&message); err != nil {
		log.Print(err)
	}
}

func TestGetAllMessage(t *testing.T) {
	messageList, err := repo.GetAllMessage()
	log.Print(messageList, err)
}
