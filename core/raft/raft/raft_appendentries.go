package raft

import (
	"log"
	"net/rpc"
	"time"
)

// each entry contains command for state machine,
// and term when entry was received by leader (first index is 1)
type LogEntry struct {
	Index   int
	Command interface{}
	Term    int
}

type AppendEntriesArgs struct {
	Term         int // leader’s term
	LeaderId     int // so follower can redirect clients
	PrevLogIndex int // index of log entry immediately preceding new ones
	PrevLogTerm  int // term of prevLogIndex entry
	// log entries to store (empty for heartbeat;
	// may send more than one for efficiency)
	Entries      []LogEntry
	LeaderCommit int // leader’s commitIndex
}

type AppendEntriesReply struct {
	Term          int  // currentTerm, for leader to update itself
	Success       bool // true if follower contained entry matching prevLogIndex and prevLogTerm
	ConflictIndex int
	ConflictTerm  int
}

func (rf *Raft) heartBeatTimer() {
	defer log.Printf("exit HeartBeatTimer: %v ", rf.me)
	appendEntriesTime := 100
	duration := time.Duration(appendEntriesTime) * time.Millisecond
	for !rf.killed() {
		time.Sleep(duration)
		rf.mu.Lock()
		if rf.state != LEADER {
			rf.mu.Unlock()
			continue
		}
		if rf.resetHeartBeatTimer {
			rf.resetHeartBeatTimer = false
			rf.mu.Unlock()
			continue
		}
		term := rf.currentTerm
		rf.mu.Unlock()
		rf.sendAppendEntries(-1, term, "heartbeat")
	}
}

// send AppendEntries RPC to peer
// set target -1 to broadcast
func (rf *Raft) sendAppendEntries(target int, term int, why string) {
	for i, p := range rf.peers {
		if i == rf.me || target != -1 && i != target {
			continue
		}
		go func(i int, p *rpc.Client) {
			rf.mu.Lock()
			if rf.isOutdate(LEADER, term) {
				rf.mu.Unlock()
				return
			}
			prevIndex := rf.nextIndex[i] - 1
			if prevIndex < rf.log[0].Index {
				installSnapshotArgs := &InstallSnapshotArgs{
					Term:              term,
					LastIncludedIndex: rf.log[0].Index,
					LastIncludedTerm:  rf.log[0].Term,
					Snapshot:          rf.snapshot,
				}
				log.Printf("%v Sending InstallSnapshot %v<%v", rf.me, prevIndex, rf.log[0].Index)
				rf.mu.Unlock()
				rf.sendInstallSnapshot(i, term, installSnapshotArgs)
				return
			}
			reply := &AppendEntriesReply{}
			args := &AppendEntriesArgs{
				Term:         rf.currentTerm,
				LeaderId:     rf.me,
				PrevLogIndex: prevIndex,
				PrevLogTerm:  rf.getLog(prevIndex).Term,
				LeaderCommit: rf.commitIndex,
			}
			if rf.nextIndex[i] <= rf.lastLogIndex() {
				args.Entries = append(make([]LogEntry, 0), rf.log[(rf.nextIndex[i]-rf.log[0].Index):]...)
			}
			// log.Printf("%v appendEntires(%v) %v => %v", rf.me, why, i, args)
			rf.mu.Unlock()
			err := p.Call("Raft.AppendEntries", args, reply)
			ok := err == nil
			if ok {
				rf.mu.Lock()
				defer rf.mu.Unlock()
				// ignore outdated AppendEntries
				if rf.isOutdate(LEADER, term) {
					return
				}
				if !reply.Success {
					if reply.ConflictTerm == -1 {
						if reply.Term > rf.currentTerm {
							rf.updateTerm(reply.Term)
							return
						}
					} else if reply.ConflictTerm > 0 {
						// If AppendEntries fails because of log inconsistency:
						// decrement nextIndex and retry (§5.3)
						// if reply.ConflictIndex<rf.lastLogIndex() {
						// 	rf.nextIndex[i] = rf.lastLogIndex()
						// }
						foundTerm := false
						for j := reply.ConflictIndex; j > rf.lastLogIndex(); j-- {
							if rf.getLog(j).Term == reply.ConflictTerm {
								rf.nextIndex[i] = j + 1
								foundTerm = true
								break
							} else if rf.getLog(j).Term < reply.ConflictTerm {
								break
							}
						}
						if !foundTerm {
							rf.nextIndex[i] = reply.ConflictIndex
						}
					} else {
						if reply.ConflictIndex != 0 {
							rf.nextIndex[i] = reply.ConflictIndex // set nextIndex = conflictIndex.
						} else {
							rf.nextIndex[i] = 1
						}
					}
					log.Printf("%v(term %v) decrementing nextIndex for %v to %v <=%v\n", rf.me, rf.currentTerm, i, rf.nextIndex[i], reply)
					rf.resetHeartBeatTimer = true
					go rf.sendAppendEntries(i, term, "resend")
					return
				}
				// If successful: update nextIndex and matchIndex for
				// follower (§5.3)
				rf.matchIndex[i] = prevIndex + len(args.Entries)
				rf.nextIndex[i] = rf.matchIndex[i] + 1
				rf.updateCommitIndex()
			}
		}(i, p)
	}
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	reply.Term = rf.currentTerm
	reply.Success = false

	// reply false if term < currentTerm(5.1)
	if rf.currentTerm > args.Term {
		reply.ConflictTerm = -1
		log.Printf("%v rejected: term: %v > %v\n", rf.me, rf.currentTerm, args.Term)
		return
	}

	rf.resetElectionTimer = true
	rf.leader = args.LeaderId
	// If RPC request or response contains term T > currentTerm:
	// set currentTerm = T, convert to follower (§5.1)
	if rf.currentTerm < args.Term {
		rf.updateTerm(args.Term)
	}

	// Reply false if log doesn’t contain an entry at prevLogIndex whose
	// term matches prevLogTerm (§5.3)
	if rf.lastLogIndex() < args.PrevLogIndex { // log doesn’t contain an entry at prevLogIndex
		reply.ConflictIndex = rf.lastLogIndex()
		if rf.lastLogIndex() != 0 && rf.getLog(rf.lastLogIndex()).Term != args.PrevLogTerm {
			reply.ConflictTerm = rf.getLog(rf.lastLogIndex()).Term
			log.Printf("%v rejected: inconsistent prevlog(as conflict): %v != %v => %v\n", rf.me, rf.getLog(rf.lastLogIndex()).Term, args.PrevLogTerm, reply)
		} else {
			log.Printf("%v rejected: doesn’t contain prevlog: => %v\n", rf.me, reply)
		}
		return
	}
	if args.PrevLogIndex < rf.log[0].Index {
		reply.ConflictIndex = rf.log[0].Index
		return
	}
	if rf.getLog(args.PrevLogIndex).Term != args.PrevLogTerm { // log on prevLogIndex doesn't matches prevLogTerm
		reply.ConflictTerm = rf.getLog(args.PrevLogIndex).Term
		for i := args.PrevLogIndex; i >= 0; i-- {
			if rf.getLog(i).Term != reply.ConflictTerm {
				reply.ConflictIndex = i + 1
				break
			}
		}
		log.Printf("%v rejected: inconsistent prevlog: %v != %v => %v\n", rf.me, rf.getLog(args.PrevLogIndex).Term, args.PrevLogTerm, reply)
		return
	}

	// If an existing entry conflicts with a new one (same index
	// but different terms), delete the existing entry and all that
	// follow it (§5.3)
	for _, log := range args.Entries {
		if log.Index <= rf.lastLogIndex() {
			if rf.getLog(log.Index).Term != log.Term {
				rf.log = rf.log[:(log.Index - rf.log[0].Index)]
				rf.log = append(rf.log, log)
			}
		} else {
			rf.log = append(rf.log, log)
		}
	}
	rf.persist()
	//
	// If leaderCommit > commitIndex, set commitIndex =
	// min(leaderCommit, index of last new entry)
	rf.commitIndex = min(args.LeaderCommit, rf.lastLogIndex())

	//
	// no longer used
	// now use a routine for log commitment
	//
	// Once a follower learns that a log entry is committed,
	// it applies the entry to its local state machine (in log order).
	// go rf.commitLog()

	reply.Success = true
	reply.Term = rf.currentTerm
}
