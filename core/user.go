package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/templari/shire-client/core/model"
	"github.com/templari/shire-client/core/util"
)

type CreateTokenRequest struct {
	Password string `json:"password"`
}

func MakeHttpRequest[T any](method string, url string, data interface{}, returnValue T) error {
	var body io.Reader = nil
	if data != nil {
		bytesData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		body = bytes.NewReader(bytesData)
	}
	req, _ := http.NewRequest(method, url, body)
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	responseBody, _ := io.ReadAll(resp.Body)
	if err = json.Unmarshal(responseBody, &returnValue); err != nil {
		return err
	}
	return nil
}

func (c *Core) Login(id int, password string) {
	client := &http.Client{}
	bytesData, err := json.Marshal(CreateTokenRequest{
		Password: password,
	})
	if err != nil {
		log.Print(err)
	}
	requestURL := fmt.Sprintf("%v/users/%v/token", c.InfoServerAddress, id)
	req, _ := http.NewRequest("POST", requestURL, bytes.NewReader(bytesData))
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)

	rep := model.LoginResponse{}
	// requestURL := fmt.Sprintf("%v/users/%v/token", c.InfoServerAddress, id)
	// if err := MakeHttpRequest("POST", requestURL, CreateTokenRequest{Password: password}, &rep); err != nil {
	// 	log.Fatal(err)
	// }
	if err = json.Unmarshal(body, &rep); err != nil {
		log.Fatal(err)
	}
	c.user = rep.User
	c.token = rep.Token
	log.Println(rep)

	listener, _ := util.CreateListener()
	go StartHttpServer(c, listener)
	res := strings.Split(listener.Addr().String(), ":")
	port, _ := strconv.Atoi(res[len(res)-1])
	c.UpdateUser(port)
}

func (c *Core) UpdateUser(port int) {
	client := &http.Client{}
	address, _ := util.GetIP()
	bytesData, _ := json.Marshal(model.UpdateUserRequest{
		Address: address,
		Port:    port,
	})
	requestURL := fmt.Sprintf("%v/users/%v", c.InfoServerAddress, c.user.Id)
	req, _ := http.NewRequest("PUT", requestURL, bytes.NewReader(bytesData))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", c.token)
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
	}
	rep := model.User{}
	body, _ := io.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &rep); err != nil {
		log.Print(err)
	}
	c.user = rep
}
