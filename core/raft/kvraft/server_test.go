package kvraft_test

import (
	"fmt"
	"testing"

	"github.com/templari/shire-client/core/raft/kvraft"
)

func TestKVServer(t *testing.T) {

	config := kvraft.ServerConfig{
		Endnames: make([]string, 3),
		Me:       1,
	}
	for i := 0; i < 3; i++ {
		config.Endnames[i] = fmt.Sprint("localhost:800", i)
	}
	kv, err := config.StartServer()
	if err != nil {
		t.Error(err)
	}
	client:=
	// kv.PutAppend("a", "b")

}
