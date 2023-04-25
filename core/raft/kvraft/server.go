package kvraft

import (
	"fmt"
	"log"
	"net/rpc"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/templari/shire-client/core/raft/labgob"
	"github.com/templari/shire-client/core/raft/raft"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type OpType string

const (
	GET    = "GET"
	PUT    = "PUT"
	APPEND = "APPEND"
)

type Op struct {
	Identifier
	Type  OpType
	Key   string
	Value string
}

func (op *Op) SetType(t string) {
	if t == "Get" {
		op.Type = GET
	} else if t == "Put" {
		op.Type = PUT
	} else {
		op.Type = APPEND
	}
}

type OpResult struct {
	Value string
	Err   Err
}

type historyOpResult struct {
	OpResult
	Identifier
}

type KVServer struct {
	mu      sync.Mutex
	me      int
	rf      *raft.Raft
	applyCh chan raft.ApplyMsg
	watchCh chan interface{}

	dead int32 // set by Kill()

	maxraftstate int // snapshot if log grows this big

	// Your definitions here.
	persister  *raft.Persister
	kvMap      map[string]string
	chanMap    map[int]chan OpResult
	historyMap map[uuid.UUID]historyOpResult
}

func (kv *KVServer) execute(op Op) (res OpResult) {
	log.Printf("server %v exec %v %v %v", kv.me, op.Type, op.Key, op.Value)
	v, ok := kv.kvMap[op.Key]
	res.Err = OK
	res.Value = ""
	if op.Type == APPEND {
		if !ok {
			kv.kvMap[op.Key] = op.Value
		} else {
			kv.kvMap[op.Key] = fmt.Sprintf("%v|%v", v, op.Value)
		}
	} else if op.Type == PUT {
		kv.kvMap[op.Key] = op.Value
	} else {
		if !ok {
			res.Err = ErrNoKey
		}
		res.Value = v
	}
	return
}
func (kv *KVServer) Get(args *GetArgs, reply *GetReply) (error error) {
	// log.Printf("%v Get %v", kv.me, args.Key)
	v, err := kv.callStart(&Op{
		Identifier: args.Identifier,
		Type:       GET,
		Key:        args.Key,
	})
	reply.Leader = -1
	reply.Err = err
	reply.Value = v
	if err == ErrWrongLeader {
		reply.Leader = kv.rf.GetLeader()
	}
	return nil
}

func (kv *KVServer) PutAppend(args *PutAppendArgs, reply *PutAppendReply) (error error) {
	// log.Printf("%v %v key:%v value:%v", kv.me, args.Op, args.Key, args.Value)
	op := &Op{
		Identifier: args.Identifier,
		Key:        args.Key,
		Value:      args.Value,
	}
	op.SetType(args.Op)
	_, err := kv.callStart(op)
	reply.Err = err
	if err == ErrWrongLeader {
		reply.Leader = kv.rf.GetLeader()
	}
	return nil
}

func (kv *KVServer) callStart(op *Op) (string, Err) {
	kv.mu.Lock()
	historyResult, ok := kv.historyMap[op.Identifier.Id]
	kv.mu.Unlock()
	if ok {
		if historyResult.Sequence == op.Sequence {
			log.Printf("%v called(cache) %v", op, historyResult)
			return historyResult.Value, historyResult.Err
		}
		// if historyResult.Sequence > op.Sequence {
		// 	return "", ErrWrongLeader
		// }
	}
	index, _, leader := kv.rf.Start(*op)
	if !leader {
		return "", ErrWrongLeader
	}
	log.Printf("%v calling %v", kv.me, op)
	opResultChan := make(chan OpResult)
	kv.mu.Lock()
	kv.chanMap[index] = opResultChan
	kv.mu.Unlock()
	select {
	case res := <-opResultChan:
		log.Printf("%v called %v", op, res)
		return res.Value, res.Err
	case <-time.After(500 * time.Millisecond):
		log.Printf("%v timeout call %v", kv.me, op)
		return "", ErrWrongLeader
	}
}

// listen for applyCh
func (kv *KVServer) applyChHandler() {
	for !kv.killed() {
		msg := <-kv.applyCh
		log.Println(kv.me, "applyChHandler", msg)
		// CommandValid==false for snapshot
		if !msg.CommandValid {
			kv.readPersist(msg.Snapshot)
			continue
		}
		op, _ := msg.Command.(Op)
		kv.mu.Lock()
		var res OpResult
		// check if duplicated
		if previousResult, ok := kv.historyMap[op.Id]; ok && previousResult.Sequence == op.Sequence {
			res = previousResult.OpResult
		} else {
			res = kv.execute(op)
			kv.historyMap[op.Id] = historyOpResult{
				res,
				op.Identifier,
			}
		}
		size := kv.persister.RaftStateSize()
		if size > kv.maxraftstate && kv.maxraftstate != -1 {
			snapshot := kv.createSnapshot()
			log.Printf("%v create snapshot %v>%v", kv.me, size, kv.maxraftstate)
			kv.rf.Snapshot(msg.CommandIndex, snapshot)
		}

		// notify watcher
		key := op.Key
		value := kv.kvMap[key]
		kv.watchCh <- value

		resCh, ok := kv.chanMap[msg.CommandIndex]
		if !ok {
			kv.mu.Unlock()
			continue
		}
		delete(kv.chanMap, msg.CommandIndex)
		kv.mu.Unlock()

		resCh <- res
	}
}

// the tester calls Kill() when a KVServer instance won't
// be needed again. for your convenience, we supply
// code to set rf.dead (without needing a lock),
// and a killed() method to test rf.dead in
// long-running loops. you can also add your own
// code to Kill(). you're not required to do anything
// about this, but it may be convenient (for example)
// to suppress debug output from a Kill()ed instance.
func (kv *KVServer) Kill() {
	atomic.StoreInt32(&kv.dead, 1)
	kv.rf.Kill()
	// Your code here, if desired.
}

func (kv *KVServer) killed() bool {
	z := atomic.LoadInt32(&kv.dead)
	return z == 1
}

func (kv *KVServer) GetRaft() *raft.Raft {
	return kv.rf
}

// servers[] contains the ports of the set of
// servers that will cooperate via Raft to
// form the fault-tolerant key/value service.
// me is the index of the current server in servers[].
// the k/v server should store snapshots through the underlying Raft
// implementation, which should call persister.SaveStateAndSnapshot() to
// atomically save the Raft state along with the snapshot.
// the k/v server should snapshot when Raft's saved state exceeds maxraftstate bytes,
// in order to allow Raft to garbage-collect its log. if maxraftstate is -1,
// you don't need to snapshot.
// StartKVServer() must return quickly, so it should start goroutines
// for any long-running work.
func StartKVServer(servers []*rpc.Client, me int, persister *raft.Persister, maxraftstate int, watchCh chan interface{}) *KVServer {
	// call labgob.Register on structures you want
	// Go's RPC library to marshall/unmarshall.
	labgob.Register(Op{})

	kv := &KVServer{
		me:           me,
		maxraftstate: maxraftstate,
		persister:    persister,
		applyCh:      make(chan raft.ApplyMsg),
		watchCh:      watchCh,
		chanMap:      map[int]chan OpResult{},
		kvMap:        make(map[string]string),
		historyMap:   make(map[uuid.UUID]historyOpResult),
	}
	log.Printf("%v server created: maxraftstate:%v", me, maxraftstate)
	kv.rf = raft.Make(servers, me, persister, kv.applyCh)

	// kv.readPersist(persister.ReadSnapshot())
	go kv.applyChHandler()
	return kv
}
