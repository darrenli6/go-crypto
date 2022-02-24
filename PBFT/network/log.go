package network

import (
	"fmt"

	"goland/go-crypto/PBFT/consensus"
)

func LogMsg(msg interface{}) {

	switch msg.(type) {

	case *consensus.RequestMsg:
		reqMsg := msg.(*consensus.RequestMsg)
		fmt.Printf("[Request] clientID :%s , Timestamp: %d,Operation :%s \n", reqMsg.ClientID, reqMsg.Timestamp, reqMsg.Operation)

	case *consensus.PrePrepareMsg:
		prePrepareMsg := msg.(*consensus.PrePrepareMsg)
		fmt.Printf("[PREPREPARE] clientID :%s ,Operation :%s,SequenceID:%d \n", prePrepareMsg.RequestMsg.ClientID, prePrepareMsg.RequestMsg.Operation, prePrepareMsg.SequenceID)
	case *consensus.VoteMsg:
		voteMsg := msg.(*consensus.VoteMsg)
		if voteMsg.MsgType == consensus.PrepareMsg {
			fmt.Printf("[prepare] nodeid %s \n", voteMsg.NodeID)
		} else if voteMsg.MsgType == consensus.CommitMsg {
			fmt.Printf(" [commit] NodeId : %s \n", voteMsg.NodeID)
		}
	}

}

func LogStage(stage string, isDone bool) {

	if isDone {
		fmt.Printf("[stage-done] %s \n", stage)
	} else {
		fmt.Printf("[stage-begin] %s \n", stage)
	}

}
