package kvraft

import (
	"bytes"
	"log"

	"github.com/google/uuid"
	"github.com/templari/shire-client/core/raft/labgob"
)

type Snapshot struct {
	HistoryMap map[uuid.UUID]historyOpResult
	RaftState  []byte
}

func (kv *KVServer) createSnapshot() []byte {
	w := new(bytes.Buffer)
	e := labgob.NewEncoder(w)
	e.Encode(kv.kvMap)
	e.Encode(kv.historyMap)
	data := w.Bytes()
	return data
}

func (kv *KVServer) readPersist(snapshot []byte) {
	if len(snapshot) == 0 {
		return
	}
	r := bytes.NewBuffer(snapshot)
	d := labgob.NewDecoder(r)
	kvMap := make(map[string]string)
	historyMap := make(map[uuid.UUID]historyOpResult)
	if d.Decode(&kvMap) != nil || d.Decode(&historyMap) != nil {
		log.Printf("%v error reading persistence ", kv.me)
		return
	} else {
		kv.kvMap = kvMap
		kv.historyMap = historyMap
		log.Printf("%v reading persistence %v", kv.me, kv.historyMap)
	}
}
