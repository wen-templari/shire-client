package raft

import (
	"log"
	"time"
)

// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in part 2D you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh, but set CommandValid to false for these
// other uses.
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int

	// For 2D:
	SnapshotValid bool
	Snapshot      []byte
	SnapshotTerm  int
	SnapshotIndex int
}

func (rf *Raft) checkCommitTimer() {
	defer log.Printf("exit CheckCommitTimer: %v ", rf.me)
	checkCommitTime := 30
	duration := time.Duration(checkCommitTime) * time.Millisecond
	for !rf.killed() {
		time.Sleep(duration)
		rf.commitLog()
	}
}

// If there exists an N such that N > commitIndex, a majority
// of matchIndex[i] ≥ N, and log[N].term == currentTerm:
// set commitIndex = N (§5.3, §5.4).
func (rf *Raft) updateCommitIndex() {
	majority := len(rf.peers)/2 + 1
	if rf.commitIndex == 0 && rf.log[0].Index != 0 {
		rf.commitIndex = rf.log[0].Index
	}
	for n := rf.commitIndex + 1; n <= rf.lastLogIndex(); n++ {
		if rf.getLog(n).Term != rf.currentTerm {
			continue
		}
		counter := 1
		for _, index := range rf.matchIndex {
			if index >= n {
				counter++
			}
		}
		if counter < majority {
			return
		}
		rf.commitIndex = n
		// log.Printf("setting commitIndex : %v |  %v/%v\n", rf.commitIndex, counter, len(rf.peers))
	}
}

// If commitIndex > lastApplied: increment lastApplied,
// apply log[lastApplied] to state machine (§5.3)
//
// apply logs from lastApplied+1 to commitIndex(including)
func (rf *Raft) commitLog() {
	rf.mu.Lock()
	if rf.lastApplied >= rf.commitIndex {
		rf.mu.Unlock()
		return
	}
	entries := append([]LogEntry{}, rf.log[rf.lastApplied+1-rf.log[0].Index:rf.commitIndex+1-rf.log[0].Index]...)
	rf.mu.Unlock()
	for _, entry := range entries {
		applyMsg := ApplyMsg{
			CommandValid: true,
			Command:      entry.Command,
			CommandIndex: entry.Index,
		}
		log.Printf("log %v (%v) committed by %v \n", applyMsg.CommandIndex, applyMsg.Command, rf.me)
		rf.applyCh <- applyMsg
		rf.mu.Lock()
		rf.lastApplied = entry.Index
		rf.mu.Unlock()
	}
}
