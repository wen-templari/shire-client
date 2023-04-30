package core

import (
	"log"
	"net"
	"net/http"

	"github.com/templari/shire-client/core/util"
	"github.com/templari/shire-client/model"
)

type Core struct {
	InfoServerAddress string
	user              model.User
	token             string
	subscribers       []chan model.Message
	listener          net.Listener

	rpcMux      *http.ServeMux
	rpcListener net.Listener

	wrappers map[int]*RaftWrapper
}

func (c *Core) Logout() {
	c.user = model.User{}
	c.token = ""
	c.subscribers = make([]chan model.Message, 0)
	c.wrappers = make(map[int]*RaftWrapper)
}

func MakeCore(address string) *Core {
	return &Core{
		InfoServerAddress: address,
		subscribers:       make([]chan model.Message, 0),
		wrappers:          make(map[int]*RaftWrapper),
	}
}

func (c *Core) startServer() {
	c.listener, _ = util.CreateListener()
	c.rpcListener, _ = util.CreateListener()
	c.rpcMux = http.NewServeMux()
	go http.Serve(c.rpcListener, c.rpcMux)
	log.Printf("Login as %v,listening at %v", c.user.Name, c.listener.Addr().String())
	go StartHttpServer(c, c.listener)

	// res := strings.Split(c.listener.Addr().String(), ":")
	// port, _ := strconv.Atoi(res[len(res)-1])
	c.UpdateUser()
}

func (c *Core) Ping() string {
	return "pong"
}
