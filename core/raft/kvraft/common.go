package kvraft

import "github.com/google/uuid"

const (
	OK             = "OK"
	ErrNoKey       = "ErrNoKey"
	ErrWrongLeader = "ErrWrongLeader"
)

type Err string

type Identifier struct {
	Id       uuid.UUID
	Sequence int
}

// Put or Append
type PutAppendArgs struct {
	Identifier
	Key   string
	Value string
	Op    string // "Put" or "Append"
}

type PutAppendReply struct {
	Err    Err
	Leader int
}

type GetArgs struct {
	Identifier
	Key string
}

type GetReply struct {
	Err    Err
	Leader int
	Value  string
}
