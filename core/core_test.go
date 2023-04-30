package core_test

import (
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

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

type Host struct {
	id       int
	name     string
	password string
	replies  []string
	speed    int
	mu       sync.Mutex
	bucket   int
}

func TestStartHostGroup(t *testing.T) {
	hostGroup := make([]*Host, 0)
	mando := &Host{
		name:     "mando",
		id:       2,
		password: "12345",
		replies: []string{
			"This is the way.", "I have spoken",
			"I can bring you in warm, or I can bring you in cold.",
			"I like those odds.",
			"… Bounty hunting is a complicated profession.",
			"I’m a Mandalorian. Weapons are part of my religion.",
			"Stop touching things.",
			"Wherever I go, he goes",
			"Dank farrik.",
			"I’ll see you again. I promise.",
		},
		speed: 1,
	}
	hostGroup = append(hostGroup, mando)
	vader := &Host{
		name:     "vader",
		id:       3,
		password: "12345",
		replies: []string{
			"He’s as clumsy as he is stupid.",
			"You don’t know the power of the dark side!",
			"It will soon see the end of the rebellion.",
			"I am altering the deal. Pray I don’t alter it any further.",
			"Be careful not to choke on your aspirations.",
			"You have controlled your fear. Now, release your anger. only your hatred can destroy me.",
			"When I left you, I was but the learner. Now I am the master.",
			"Don’t be too proud of this technological terror you’ve constructed. The ability to destroy a planet is insignificant next to the power of the force.",
			"I find your lack of faith disturbing.",
			"Anakin Skywalker was weak. I destroyed him.",
			"Apology accepted.",
		},
		speed: 1,
	}
	hostGroup = append(hostGroup, vader)
	for _, v := range hostGroup {
		v.start()
	}

	for {

	}
}

func (h *Host) start() error {
	c := core.MakeCore(infoServerAddr)

	if _, err := c.Login(h.id, h.password); err != nil {
		return err
	}
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				h.mu.Lock()
				if h.bucket < 3 {
					h.bucket++
				}
				h.mu.Unlock()
				// do stuff
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	messageChan := make(chan model.Message, 10)
	c.Subscribe(messageChan)
	go func(h *Host, messageChan chan model.Message) {
		for {
			message := <-messageChan
			var toWhomShouldIReply int
			// TODO when is group
			h.mu.Lock()
			if message.From == c.GetUser().Id || h.bucket < 0 {
				h.mu.Unlock()
				continue
			} else {
				h.bucket--
				h.mu.Unlock()
				toWhomShouldIReply = message.From
			}
			msg := model.Message{
				From:    c.GetUser().Id,
				To:      toWhomShouldIReply,
				GroupId: message.GroupId,
				Content: h.replies[rand.Intn(len(h.replies))],
			}
			rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
			time.Sleep(time.Duration((500 + len(msg.Content)*30) * h.speed * int(time.Millisecond)))
			err := c.SendMessage(msg)
			if err != nil {
				log.Println(err)
			}

		}
	}(h, messageChan)
	return nil
}

// func TestStartMandoServer(t *testing.T) {
// 	mando := core.MakeCore(infoServerAddr)

// 	if _, err := mando.Login(2, "12345"); err != nil {
// 		t.Error(err)
// 	}
// 	replies := []string{"This is the way.", "I have spoken",
// 		"I can bring you in warm, or I can bring you in cold.",
// 		"I like those odds.",
// 		"… Bounty hunting is a complicated profession.",
// 		"I’m a Mandalorian. Weapons are part of my religion.",
// 		"Stop touching things.",
// 		"Wherever I go, he goes",
// 		"Dank farrik.",
// 		"I’ll see you again. I promise."}
// 	messageChan := make(chan model.Message, 10)
// 	mando.Subscribe(messageChan)
// 	for {
// 		message := <-messageChan
// 		var toWhomShouldIReply int
// 		// TODO when is group
// 		if message.From == mando.GetUser().Id {
// 			continue
// 		} else {
// 			toWhomShouldIReply = message.From
// 		}
// 		msg := model.Message{
// 			From:    mando.GetUser().Id,
// 			To:      toWhomShouldIReply,
// 			Content: replies[rand.Intn(len(replies))],
// 		}
// 		// rand.Seed(time.Now().Unix()) // initialize global pseudo random generator
// 		// time.Sleep(1000 * time.Microsecond)
// 		err := mando.SendMessage(msg)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 	}

// }

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
