package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isLeader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	"log"
	"net/rpc"
	"sync"
	"sync/atomic"
)

type state int

const (
	FOLLOWER = state(iota)
	CANDIDATE
	LEADER
)

// A Go object implementing a single Raft peer.
type Raft struct {
	mu        sync.Mutex    // Lock to protect shared access to this peer's state
	peers     []*rpc.Client // RPC end points of all peers
	persister *Persister    // Object to hold this peer's persisted state
	me        int           // this peer's index into peers[]
	leader    int
	dead      int32 // set by Kill()

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

	// latest term server has seen (initialized to 0 on first boot, increases monotonically)
	currentTerm int
	votedFor    int        // candidateId that received vote in current term (or null if none)
	log         []LogEntry // log entries;
	// index of highest log entry known to be committed
	// (initialized to 0, increases monotonically)
	commitIndex int
	// index of highest log entry applied to state machine (initialized to 0, increases monotonically)
	lastApplied int
	// for each server, index of the next log entry to send to that server
	// (initialized to leader last log index + 1)
	nextIndex []int
	// for each server, index of highest log entry known to be replicated on server
	// (initialized to 0, increases monotonically)
	matchIndex []int

	state    state
	snapshot []byte

	resetElectionTimer  bool
	resetHeartBeatTimer bool

	applyCh chan ApplyMsg
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {
	var term int
	var isLeader bool
	// Your code here (2A).
	rf.mu.Lock()
	defer rf.mu.Unlock()
	term = rf.currentTerm
	isLeader = rf.state == LEADER
	return term, isLeader
}

// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	// Your code here (2B).
	index := -1
	isLeader := true

	rf.mu.Lock()
	defer rf.mu.Unlock()
	term := rf.currentTerm
	if rf.state != LEADER {
		isLeader = false
		return index, term, isLeader
	}
	// append log
	index = rf.nextIndex[rf.me]
	rf.nextIndex[rf.me]++
	logToAppend := &LogEntry{
		Index:   index,
		Command: command,
		Term:    rf.currentTerm,
	}
	rf.log = append(rf.log, *logToAppend)
	log.Printf("%v appending %v", rf.me, logToAppend)
	rf.persist()

	// issues AppendEntries RPCs in parallel to each of the other servers to replicate the entry.
	// this seems to be necessary to reach the speed required by the test (kvraft GenericTestSpeed).
	rf.resetHeartBeatTimer = true
	go rf.sendAppendEntries(-1, term, "replicate")

	return index, term, isLeader
}

// the tester doesn't halt goroutines created by Raft after each test,
// but it does call the Kill() method. your code can use killed() to
// check whether Kill() has been called. the use of atomic avoids the
// need for a lock.
//
// the issue is that long-running goroutines use memory and may chew
// up CPU time, perhaps causing later tests to fail and generating
// confusing debug output. any goroutine with a long-running loop
// should call killed() to check whether it should stop.
func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
	// Your code here, if desired.
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}

// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
func Make(peers []*rpc.Client, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{
		peers:               peers,
		persister:           persister,
		leader:              0,
		me:                  me,
		commitIndex:         0,
		currentTerm:         0,
		lastApplied:         0,
		votedFor:            -1,
		state:               FOLLOWER,
		log:                 make([]LogEntry, 0),
		nextIndex:           make([]int, len(peers)),
		matchIndex:          make([]int, len(peers)),
		resetElectionTimer:  false,
		resetHeartBeatTimer: false,
		applyCh:             applyCh,
		snapshot:            nil,
	}
	// log.SetOutput(ioutil.Discard)
	rf.mu.Lock()
	defer rf.mu.Unlock()
	rf.log = append(rf.log, LogEntry{
		Term:  0,
		Index: 0,
	})
	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())
	rf.snapshot = persister.ReadSnapshot()
	rf.lastApplied = rf.log[0].Index
	rf.printStatus("created")

	go rf.heartBeatTimer()
	go rf.electionTimer()
	go rf.checkCommitTimer()

	go rf.statusReport()
	return rf
}
