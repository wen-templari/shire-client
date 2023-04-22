package raft

import (
	"log"
	"time"
)

// Debugging
const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

func (rf *Raft) GetLeader() int {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	return rf.leader
}

// If RPC request or response contains term T > currentTerm:
// set currentTerm = T, convert to follower (§5.1)
func (rf *Raft) updateTerm(term int) {
	// rf.state = FOLLOWER
	rf.currentTerm = term
	rf.votedFor = -1
	if rf.state != FOLLOWER {
		rf.state = FOLLOWER
		log.Printf("step down: %v \n ", rf.me)
	}
	rf.persist()
}

// return true is is outdated
// pass the state and term you except raft to be
func (rf *Raft) isOutdate(state state, term int) bool {
	if state == rf.state && term == rf.currentTerm {
		return false
	} else {
		return true
	}
}

func (rf *Raft) getLog(index int) LogEntry {
	translatedIndex := index - rf.log[0].Index
	if translatedIndex < 0 || translatedIndex > len(rf.log)-1 {
		log.Printf("%v getLog index error: %v - %v", rf.me, index, rf.log[0].Index)
	}
	return rf.log[index-rf.log[0].Index]
}

func (rf *Raft) lastLogIndex() int {
	index := rf.log[len(rf.log)-1].Index
	return index
}

func (rf *Raft) statusReport() {
	defer log.Printf("exit statusReport: %v", rf.me)
	for !rf.killed() {
		time.Sleep(time.Millisecond * 400)
		rf.mu.Lock()
		rf.printStatus("report")
		rf.mu.Unlock()
	}
}

func (rf *Raft) printStatus(subject string) {
	statusList := [3]string{"FOLLOWER", "CANDIDATE", "LEADER"}
	log.Printf("%v: %v | %v | Term %v | CommitIndex %v | %v \n", subject, rf.me, statusList[rf.state], rf.currentTerm, rf.commitIndex, rf.log)
}

func min(a int, b int) int {
	if a >= b {
		return b
	} else {
		return a
	}
}
