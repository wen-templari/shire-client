package kvraft

import (
	"log"
	"math/rand"
	"net/rpc"
	"sync"

	"github.com/google/uuid"
)

type Clerk struct {
	servers []*rpc.Client
	// You will have to modify this struct.
	id           uuid.UUID
	sequence     int
	recentLeader int
	mu           sync.Mutex
}

// func nrand() int64 {
// 	max := big.NewInt(int64(1) << 62)
// 	bigx, _ := rand.Int(rand.Reader, max)
// 	x := bigx.Int64()
// 	return x
// }

func MakeClerk(servers []*rpc.Client) *Clerk {
	ck := &Clerk{
		servers:      servers,
		id:           uuid.New(),
		sequence:     0,
		recentLeader: 0,
	}
	// You'll have to add code here.
	return ck
}

// fetch the current value for a key.
// returns "" if the key does not exist.
// keeps trying forever in the face of all other errors.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("KVServer.Get", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
func (ck *Clerk) Get(key string) string {
	log.Println("get", key)
	ck.sequence++
	args := GetArgs{
		Identifier: Identifier{
			Id:       ck.id,
			Sequence: ck.sequence,
		},
		Key: key,
	}
	for {
		ck.mu.Lock()
		nextStart := ck.recentLeader
		// log.Println("starting :", nextStart)
		ck.mu.Unlock()
		for i := nextStart; i < len(ck.servers); i++ {
			reply := GetReply{}
			if err := ck.servers[i].Call("KVServer.Get", &args, &reply); err != nil {
				// ck.mu.Lock()
				// ck.setRecentLeader(0)
				// ck.mu.Unlock()
				continue
			}
			switch reply.Err {
			case ErrWrongLeader:
				ck.mu.Lock()
				// log.Printf("From %v: ErrWrongLeader %v | current %v", i, reply.Leader, ck.recentLeader)
				ck.setRecentLeader(reply.Leader)
				ck.mu.Unlock()
				continue
			case OK:
				log.Println("get reply", reply.Value)
				return reply.Value
			case ErrNoKey:
				return ""
			}
		}
	}
}

// shared by Put and Append.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("KVServer.PutAppend", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
func (ck *Clerk) PutAppend(key string, value string, op string) {
	log.Println(op, key, value)
	ck.sequence++
	args := PutAppendArgs{
		Identifier: Identifier{
			Id:       ck.id,
			Sequence: ck.sequence,
		},
		Key:   key,
		Value: value,
		Op:    op,
	}
	for {
		ck.mu.Lock()
		nextStart := ck.recentLeader
		// log.Println("starting :", nextStart)
		ck.mu.Unlock()
		for i := nextStart; i < len(ck.servers); i++ {
			// log.Println(ck.servers[i])
			reply := PutAppendReply{}
			if err := ck.servers[i].Call("KVServer.PutAppend", &args, &reply); err != nil {
				// ck.mu.Lock()
				// ck.setRecentLeader(0)
				// ck.mu.Unlock()
				continue
			}
			switch reply.Err {
			case ErrWrongLeader:
				ck.mu.Lock()
				// log.Printf("From %v: ErrWrongLeader %v | current %v", i, reply.Leader, ck.recentLeader)
				ck.setRecentLeader(reply.Leader)
				ck.mu.Unlock()
				continue
			case OK:
				log.Println("Put Append reply", reply)
				return
			}
		}
	}
}

// func (ck *Clerk) setRecentLeader(leader int) {
// 	newLeader := 0
// 	if newLeader == leader {
// 		for newLeader == leader {
// 			newLeader = rand.Intn(len(ck.servers) - 1)
// 		}
// 	} else {
// 		newLeader = leader
// 	}
// 	ck.recentLeader = newLeader
// }

func (ck *Clerk) setRecentLeader(leader int) {
	newLeader := 0
	if newLeader == 0 && leader == 0 {
		for newLeader != ck.recentLeader {
			newLeader = rand.Intn(len(ck.servers) - 1)
		}
	} else {
		newLeader = ck.recentLeader
	}
	// log.Printf("setting recentLeader %v to %v", ck.recentLeader, newLeader)
	ck.recentLeader = newLeader
}

func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}
func (ck *Clerk) Append(key string, value string) {
	ck.PutAppend(key, value, "Append")
}
