package raft

import "log"

// A service wants to switch to snapshot.  Only do so if Raft hasn't
// have more recent info since it communicate the snapshot on applyCh.
func (rf *Raft) CondInstallSnapshot(lastIncludedTerm int, lastIncludedIndex int, snapshot []byte) bool {
	// Your code here (2D).

	return true
}

type InstallSnapshotArgs struct {
	Term              int
	LeaderId          int
	LastIncludedIndex int
	LastIncludedTerm  int
	Snapshot          []byte
}

type InstallSnapshotReply struct {
	Term int
}

func (rf *Raft) InstallSnapshot(args *InstallSnapshotArgs, reply *InstallSnapshotReply) (err error) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	term := rf.currentTerm
	reply.Term = term
	if args.Term < rf.currentTerm {
		return
	}
	log.Printf("%v InstallSnapshot: %v\n", rf.me, args.LastIncludedIndex)
	rf.trimLog(rf.lastLogIndex())
	startingLog := LogEntry{
		Term:  args.LastIncludedTerm,
		Index: args.LastIncludedIndex,
	}
	rf.log[0] = startingLog
	rf.snapshot = args.Snapshot
	applyMsg := ApplyMsg{
		SnapshotValid: true,
		Snapshot:      rf.snapshot,
		SnapshotTerm:  args.LastIncludedTerm,
		SnapshotIndex: args.LastIncludedIndex,
	}
	rf.applyCh <- applyMsg
	rf.lastApplied = args.LastIncludedIndex
	rf.commitIndex = args.LastIncludedIndex
	rf.persist()
	rf.printStatus("InstallSnapshot")

	return
}

func (rf *Raft) sendInstallSnapshot(target int, term int, args *InstallSnapshotArgs) {
	log.Printf("%v Sending InstallSnapshot => %v", rf.me, target)
	reply := &InstallSnapshotReply{}
	err := rf.peers[target].Call("Raft.InstallSnapshot", args, reply)
	ok := err == nil
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if ok {
		log.Printf("%v InstallSnapshot reply received: %v", rf.me, reply)
		if reply.Term > rf.currentTerm {
			rf.updateTerm(args.Term)
			return
		}
		rf.nextIndex[target] = rf.log[0].Index + 1
		go rf.sendAppendEntries(target, term, "snapshot done")
	}
}

// the service says it has created a snapshot that has
// all info up to and including index. this means the
// service no longer needs the log through (and including)
// that index. Raft should now trim its log as much as possible.
func (rf *Raft) Snapshot(index int, snapshot []byte) {
	// Your code here (2D).
	log.Printf("%v taking snapshot at index %v\n", rf.me, index)
	rf.mu.Lock()
	defer func() {
		rf.printStatus("Snapshot Done")
		rf.mu.Unlock()
	}()
	rf.trimLog(index)
	rf.snapshot = snapshot
	rf.persist()
}

func (rf *Raft) trimLog(index int) {
	log.Printf("%v trimming %v %v \n", rf.me, index, index-rf.log[0].Index)
	rf.log = append(make([]LogEntry, 0), rf.log[(index-rf.log[0].Index):]...)
	rf.log[0].Command = nil
}
