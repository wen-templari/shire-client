package core

import (
	"fmt"
	"net/rpc"

	"github.com/templari/shire-client/core/raft/kvraft"
)

type RaftWrapper struct {
	endnames []string
	ends     []*rpc.Client
	myname   string
	me       int
	groupId  int

	watchCh chan interface{}

	Server    *kvraft.KVServer
	rpcServer *rpc.Server
}

func (wrap *RaftWrapper) ServeRPC() error {
	wrap.rpcServer = rpc.NewServer()
	wrap.rpcServer.HandleHTTP(fmt.Sprint("/rpc/", wrap.groupId), fmt.Sprint("/debug/rpc", wrap.groupId))
	return nil
}

func (wrap *RaftWrapper) Connect() error {
	if len(wrap.endnames) == 0 {
		return nil
	}
	wrap.ends = make([]*rpc.Client, len(wrap.endnames))
	for i, end := range wrap.endnames {
		if i == wrap.me {
			continue
		}
		conn, err := rpc.DialHTTP("tcp", end)
		if err != nil {
			return err
		}
		wrap.ends[i] = conn
	}
	return nil
}

func (wrap *RaftWrapper) makeKVServer() error {
	wrap.Server = kvraft.StartKVServer(wrap.ends, wrap.me, nil, -1, wrap.watchCh)

	return nil
}

func (wrap *RaftWrapper) convertChan()  {
	for {
		data := <-wrap.watchCh
		
	}
}



func MakeRaftWrapper(endnames []string, myname string,  chan interface{}) *RaftWrapper {

	var me int
	for i, name := range endnames {
		if name == myname {
			me = i
			break
		}
	}
	return &RaftWrapper{
		endnames: endnames,
		myname:   myname,
		me:       me,
		watchCh:  watchCh,
	}
}
