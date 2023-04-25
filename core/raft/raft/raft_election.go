package raft

import (
	"log"
	"math/rand"
	"time"
)

// example ReestVote RPC arguments structure.
// field names must start with capital letters!
type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	Term         int // candidate’s term
	CandidateId  int // candidate requesting vote
	LastLogIndex int // index of candidate’s last log entry (§5.4)
	LastLogTerm  int // term of candidate’s last log entry (§5.4)
}

// example RequestVote RPC reply structure.
// field names must start with capital letters!
type RequestVoteReply struct {
	// Your data here (2A).
	Who         int  // the voter
	Term        int  // currentTerm, for candidate to update itself
	VoteGranted bool // true means candidate received vote
}

func (rf *Raft) electionTimer() {
	defer log.Printf("exit ElectionTimer: %v ", rf.me)
	randElectionTime := rand.Int63()%250 + 150
	duration := time.Duration(randElectionTime) * time.Millisecond
	for !rf.killed() {
		time.Sleep(duration)
		rf.mu.Lock()
		if rf.resetElectionTimer {
			rf.resetElectionTimer = false
			rf.mu.Unlock()
			continue
		}
		if rf.state == LEADER {
			rf.mu.Unlock()
			continue
		}
		rf.mu.Unlock()
		rf.startElection()
	}
}

func (rf *Raft) startElection() {
	rf.mu.Lock()
	rf.state = CANDIDATE // convert to candidate
	rf.currentTerm++     // increment currentTerm
	term := rf.currentTerm
	rf.printStatus("election")
	rf.votedFor = rf.me // vote for self
	rf.persist()
	voteCount := 1
	rf.resetElectionTimer = true // reset election timer
	majority := len(rf.peers)/2 + 1
	lastLogIndex := rf.lastLogIndex()
	requestVoteArgs := &RequestVoteArgs{
		Term:         rf.currentTerm,
		CandidateId:  rf.me,
		LastLogIndex: lastLogIndex,
		LastLogTerm:  rf.getLog(lastLogIndex).Term,
	}
	rf.mu.Unlock()
	for peer := range rf.peers { // Send RequestVote RPCs to all other servers
		if peer == rf.me { // ignore self
			continue
		}
		go func(peer int) {
			tempReply := RequestVoteReply{}
			// client := rf.peers[peer]
			// log.Printf("%v => voteRequest: %v | client %v | Term %v | %v \n", rf.me, peer, client, term, requestVoteArgs)
			error := rf.peers[peer].Call("Raft.RequestVote", requestVoteArgs, &tempReply)
			log.Println(error)
			ok := error == nil
			rf.mu.Lock()
			defer rf.mu.Unlock()
			// ignore outdated reply
			if !ok || rf.isOutdate(CANDIDATE, term) {
				return
			}
			log.Printf("%v <= voteReply: %v | Term %v | %v \n", rf.me, tempReply.Who, tempReply.Term, tempReply.VoteGranted)
			// receive from higher term: convert to follower
			if tempReply.Term > rf.currentTerm {
				rf.updateTerm(tempReply.Term)
				rf.resetElectionTimer = true
				return
			}
			if tempReply.VoteGranted {
				voteCount++
				// If votes received from majority of servers: become leader
				if voteCount >= majority && rf.currentTerm == term {
					rf.state = LEADER
					rf.printStatus("elected")
					rf.leader = rf.me
					// Initialize nextIndex to leader last log index + 1
					for i := range rf.nextIndex {
						rf.nextIndex[i] = rf.lastLogIndex() + 1
					}
					for i := range rf.matchIndex { // Initialize to 0
						rf.matchIndex[i] = 0
					}
					// rf.resetHeartBeatTimer = true
					go rf.sendAppendEntries(-1, term, "elected")
					return
				}
			}
		}(peer)
	}
}

func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) (err error) {
	// Your code here (2A, 2B).
	rf.mu.Lock()
	defer rf.mu.Unlock()

	reply.Who = rf.me
	reply.VoteGranted = false
	reply.Term = rf.currentTerm
	if args.Term < rf.currentTerm { // Reply false if term < currentTerm (§5.1)
		return
	}

	if args.Term > rf.currentTerm {
		rf.updateTerm(args.Term)
	}
	isGivenLogMoreUpToDate := func(givenLastLogIndex int, givenTerm int) bool {
		myLastLogIndex := rf.lastLogIndex()
		lastLog := rf.getLog(myLastLogIndex)
		if lastLog.Term == givenTerm {
			// If the logs end with the same term, then whichever log is longer is more up-to-date.
			return givenLastLogIndex >= myLastLogIndex
		} else {
			// If the logs have last entries with different terms, then the log with the later term is more up-to-date.
			return lastLog.Term < givenTerm
		}
	}
	// if votedFor is null or candidateId, and candidate’s log is at
	// least as up-to-date as receiver’s log, grant vote (§5.2, §5.4)
	if (rf.votedFor == -1 || rf.votedFor == args.CandidateId) && isGivenLogMoreUpToDate(args.LastLogIndex, args.LastLogTerm) {
		reply.VoteGranted = true
		rf.votedFor = args.CandidateId
		rf.resetElectionTimer = true
		rf.persist()
		return
	}
	return nil
}
