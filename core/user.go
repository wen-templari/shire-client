package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/templari/shire-client/core/model"
	"github.com/templari/shire-client/core/util"
)

type CreateTokenRequest struct {
	Password string `json:"password"`
}

func (c *Core) GetUser() model.User {
	return c.user
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

func (c *Core) Register(name string, password string) (model.User, error) {
	client := &http.Client{}
	bytesData, err := json.Marshal(model.User{
		Name:     name,
		Password: password,
	})
	if err != nil {
		return model.User{}, err
	}
	requestURL := fmt.Sprintf("%v/users", c.InfoServerAddress)
	req, _ := http.NewRequest("POST", requestURL, bytes.NewReader(bytesData))
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return model.User{}, err
	}
	body, _ := io.ReadAll(resp.Body)

	rep := model.LoginResponse{}
	if err = json.Unmarshal(body, &rep); err != nil {
		return model.User{}, err
	}
	c.user = rep.User
	c.token = rep.Token
	listener, _ := util.CreateListener()
	go StartHttpServer(c, listener)
	res := strings.Split(listener.Addr().String(), ":")
	port, _ := strconv.Atoi(res[len(res)-1])
	c.UpdateUser(port)
	return rep.User, nil
}

func (c *Core) Login(id int, password string) (model.User, error) {
	client := &http.Client{}
	bytesData, err := json.Marshal(CreateTokenRequest{
		Password: password,
	})
	if err != nil {
		return model.User{}, err
	}
	requestURL := fmt.Sprintf("%v/users/%v/token", c.InfoServerAddress, id)
	req, _ := http.NewRequest("POST", requestURL, bytes.NewReader(bytesData))
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return model.User{}, err
	}
	body, _ := io.ReadAll(resp.Body)

	rep := model.LoginResponse{}
	// requestURL := fmt.Sprintf("%v/users/%v/token", c.InfoServerAddress, id)
	// if err := MakeHttpRequest("POST", requestURL, CreateTokenRequest{Password: password}, &rep); err != nil {
	// 	log.Fatal(err)
	// }
	if err = json.Unmarshal(body, &rep); err != nil {
		return model.User{}, err
	}
	c.user = rep.User
	c.token = rep.Token
	listener, _ := util.CreateListener()
	go StartHttpServer(c, listener)
	res := strings.Split(listener.Addr().String(), ":")
	port, _ := strconv.Atoi(res[len(res)-1])
	c.UpdateUser(port)
	return rep.User, nil
}

func (c *Core) UpdateUser(port int) (model.User, error) {
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
		return model.User{}, err
	}
	rep := model.User{}
	body, _ := io.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &rep); err != nil {
		return model.User{}, err
	}
	c.user = rep
	return rep, nil
}

func (c *Core) GetUserById(id int) (model.User, error) {
	client := &http.Client{}
	requestURL := fmt.Sprintf("%v/users/%v", c.InfoServerAddress, id)
	req, _ := http.NewRequest("GET", requestURL, nil)
	resp, err := client.Do(req)
	if err != nil {
		return model.User{}, err
	}
	rep := model.User{}
	body, _ := io.ReadAll(resp.Body)
	if err = json.Unmarshal(body, &rep); err != nil {
		return model.User{}, err
	}
	return rep, err
}
