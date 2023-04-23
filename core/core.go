package core

import (
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/templari/shire-client/core/util"
	"github.com/templari/shire-client/model"
)

type Core struct {
	InfoServerAddress string
	user              model.User
	token             string
	subscribers       []chan model.Message
	listener          net.Listener

	// endnames   []string
	// ends       []*rpc.Client
	// myname     string
	// me         int
	// kvservers  []*kvraft.KVServer
	// rpcServers []*rpc.Server
}

func MakeCore(address string) *Core {
	return &Core{
		InfoServerAddress: address,
		subscribers:       make([]chan model.Message, 0),
	}
}

func (c *Core) startServer() {
	c.listener, _ = util.CreateListener()
	log.Printf("Login as %v,listening at %v", c.user.Name, c.listener.Addr().String())
	go StartHttpServer(c, c.listener)
	res := strings.Split(c.listener.Addr().String(), ":")
	port, _ := strconv.Atoi(res[len(res)-1])
	c.UpdateUser(port)
}

func (c *Core) Ping() string {
	return "pong"
}
