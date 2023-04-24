package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/templari/shire-client/model"
)

// start a group chat
func (c *Core) CreateGroup(idList []int) (group model.Group, err error) {

	// 1. prepare group
	client := &http.Client{}
	bytesData, _ := json.Marshal(model.CreateGroupRequest{
		Ids: idList,
	})
	requestURL := fmt.Sprintf("%v/groups", c.InfoServerAddress)
	req, _ := http.NewRequest("POST", requestURL, bytes.NewReader(bytesData))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", c.token)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	body, _ := io.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &group); err != nil {
		return
	}

	// 2. call prepare group
	wg := sync.WaitGroup{}
	for _, user := range group.Users {
		wg.Add(1)
		go func(user model.User) {
			defer wg.Done()
			err := c.CallPrepare(user, group.Id)
			if err != nil {
				log.Println(err)
			}
		}(user)
	}
	wg.Wait()
	// 3. call start
	// err = c.StartGroup(group.Id)
	// if err != nil {
	// 	log.Print(err)
	// }

	wg = sync.WaitGroup{}
	for _, user := range group.Users {
		wg.Add(1)
		go func(user model.User) {
			defer wg.Done()
			err := c.CallStart(user, group.Id)
			log.Printf("calling start for %v,err: %v", user, err)
			// if err != nil {
			// 	log.Println(err)
			// }
		}(user)
	}
	wg.Wait()
	return
}

func (c *Core) CallPrepare(user model.User, groupId int) error {
	client := &http.Client{}
	requestURL := fmt.Sprintf("http://%v:%v/groups/%v/prepare", user.Address, user.Port, groupId)
	log.Println(requestURL)
	req, err := http.NewRequest("POST", requestURL, nil)
	if err != nil {
		log.Print(err)
	}
	_, err = client.Do(req)
	log.Printf("err: %v", err)
	if err != nil {
		log.Print(err)
	}
	return err
}

func (c *Core) CallStart(user model.User, groupId int) error {
	client := &http.Client{}
	requestURL := fmt.Sprintf("http://%v:%v/groups/%v/start", user.Address, user.Port, groupId)
	log.Println(requestURL)
	req, err := http.NewRequest("POST", requestURL, nil)
	if err != nil {
		log.Print(err)
	}
	_, err = client.Do(req)
	if err != nil {
		log.Print(err)
	}
	return err
}

func (c Core) GetGroupById(id int) (group model.Group, err error) {
	client := &http.Client{}
	requestURL := fmt.Sprintf("%v/groups/%v", c.InfoServerAddress, id)
	req, _ := http.NewRequest("GET", requestURL, nil)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &group)
	return
}

// start listening for rpc
func (c *Core) PrepareGroup(id int) error {
	_, ok := c.wrappers[id]
	if ok {
		// TODO
		// g.Server.Kill()
		return fmt.Errorf("group %v already exists", id)
	}
	group, err := c.GetGroupById(id)
	if err != nil {
		return err
	}
	endnames := make([]string, len(group.Users))
	me := 0
	for i, user := range group.Users {
		if user.Id == c.user.Id {
			me = i
		}
		endnames[i] = fmt.Sprintf("%v:%v", user.Address, user.RPCPort)
	}
	watchCh := make(chan interface{})
	c.wrappers[id] = MakeRaftWrapper(endnames, me, watchCh, c.rpcMux)
	if err := c.wrappers[id].ServeRPC(); err != nil {
		return err
	}
	log.Printf("user %v group %v prepared", c.user, id)
	return nil
}

func (c *Core) StartGroup(id int) error {
	w, ok := c.wrappers[id]
	if !ok {
		return fmt.Errorf("no wrapper found for group %v", id)
	}

	// set up watch channel
	go func() {
		for {
			data := <-w.watchCh
			log.Println("receive data from watch channel", data)
			// message, ok := data.(model.Message)

			// log.Printf("receive message %v,ok: %v", message, ok)
			// if ok {
			// 	c.ReceiveMessage(message)
			// }
		}
	}()

	err := w.ConnectAndStartKV()
	return err
}
