package kvraft

import (
	"net/rpc"

	"github.com/templari/shire-client/core/raft/raft"
)

type Config struct {
	endnames []string
	ends     []*rpc.Client
	myname   string
	me       int

	Server *KVServer
}

// func (cfg *Config) StartRPC() error {
// 	rpc.Register(cfg.Server)
// 	rpc.Register(cfg.Server.rf)
// 	rpc.HandleHTTP()

// 	l, e := net.Listen("tcp", ":1234")
// 	if e != nil {
// 		log.Fatal("listen error:", e)
// 	}

// 	l, e := cfg.Server.listener.Accept()
// 	if e != nil {
// 		return e
// 	}
// 	go cfg.Server.Serve(l)
// 	return nil
// }

func (cfg *Config) Connect() error {
	if len(cfg.endnames) == 0 {
		return nil
	}
	cfg.ends = make([]*rpc.Client, len(cfg.endnames))
	for i, end := range cfg.endnames {
		if i == cfg.me {
			continue
		}
		conn, err := rpc.DialHTTP("tcp", end)
		if err != nil {
			return err
		}
		cfg.ends[i] = conn
	}
	cfg.Server = StartKVServer(cfg.ends, cfg.me, raft.MakePersister(), -1)
	return nil
}

func MakeConfig(endnames []string, myname string) *Config {
	var me int
	for i, name := range endnames {
		if name == myname {
			me = i
			break
		}
	}
	return &Config{
		endnames: endnames,
		myname:   myname,
		me:       me,
	}
}
