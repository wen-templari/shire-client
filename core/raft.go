package core

import (
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"sync"

	"github.com/templari/shire-client/core/raft/kvraft"
	"github.com/templari/shire-client/core/raft/raft"
)

type RaftWrapper struct {
	endnames []string
	ends     []*rpc.Client
	me       int
	groupId  int
	mux      *http.ServeMux

	client  *kvraft.Clerk
	watchCh chan interface{}

	Server    *kvraft.KVServer
	rpcServer *rpc.Server
}

// first
func (wrap *RaftWrapper) ServeRPC() error {
	wrap.rpcServer = rpc.NewServer()
	log.Printf("me %v,group %v serving rpc", wrap.me, wrap.groupId)
	wrap.mux.Handle(fmt.Sprint("/rpc/", wrap.groupId), wrap.rpcServer)
	// wrap.rpcServer.HandleHTTP(fmt.Sprint("/rpc/", wrap.groupId), fmt.Sprint("/debug/rpc", wrap.groupId))
	return nil
}

// second
func (wrap *RaftWrapper) ConnectAndStartKV() error {
	if len(wrap.endnames) == 0 {
		return nil
	}
	wrap.ends = make([]*rpc.Client, len(wrap.endnames))
	wg := sync.WaitGroup{}
	for i, end := range wrap.endnames {
		// if i == wrap.me {
		// 	continue
		// }
		wg.Add(1)
		go func(i int, end string) {
			defer wg.Done()
			// rpcAddr := fmt.Sprintf("%v/rpc/%v", end, wrap.groupId)
			conn, err := rpc.DialHTTPPath("tcp", end, fmt.Sprint("/rpc/", wrap.groupId))
			log.Printf("me %v,group %v connecting to %v,err: %v", wrap.me, wrap.groupId, end, err)
			log.Print("conn", conn)
			// if err != nil {
			// 	return err
			// }

			wrap.ends[i] = conn
		}(i, end)
	}
	wg.Wait()
	log.Printf("me %v,group %v starting kv, conns:%v", wrap.me, wrap.groupId, wrap.ends)
	wrap.client = kvraft.MakeClerk(wrap.ends)
	wrap.Server = kvraft.StartKVServer(wrap.ends, wrap.me, raft.MakePersister(), -1, wrap.watchCh)

	wrap.rpcServer.Register(wrap.Server)
	wrap.rpcServer.Register(wrap.Server.GetRaft())
	return nil
}

func MakeRaftWrapper(endnames []string, me int, watchCh chan interface{}, mux *http.ServeMux) *RaftWrapper {
	return &RaftWrapper{
		endnames: endnames,
		me:       me,
		watchCh:  watchCh,
		mux:      mux,
	}
}
