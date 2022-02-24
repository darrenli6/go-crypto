package consensus

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// 拜占庭节点的个数
const f = 1

type State struct {
	ViewID int64

	MsgLogs *MsgLogs

	// 最后的序号id
	LastSequenceID int64

	// 当前状态
	CurrentStage Stage
}

type MsgLogs struct {
	ReqMsg      *RequestMsg
	PrepareMsgs map[string]*VoteMsg
	CommitMsgs  map[string]*VoteMsg
}

type Stage int

const (
	// 节点已经创建，共识没有开始
	Idle Stage = iota
	// 序号分配之前
	PrePrepared
	// 序号分配
	Prepared
	// 提交
	Committed
)

// lastSequenceID will be -1 if there is no last sequence ID.
// 根据视图id
func CreateState(viewID int64, lastSequenceID int64) *State {
	return &State{
		ViewID: viewID,
		MsgLogs: &MsgLogs{
			ReqMsg:      nil,
			PrepareMsgs: make(map[string]*VoteMsg),
			CommitMsgs:  make(map[string]*VoteMsg),
		},
		LastSequenceID: lastSequenceID,
		CurrentStage:   Idle,
	}
}

// 开始共识
func (state *State) StartConsensus(request *RequestMsg) (*PrePrepareMsg, error) {

	sequenceID := time.Now().UnixNano()

	//根据最后提交的信息的编号为新的请求分配编号
	if state.LastSequenceID != -1 {
		for state.LastSequenceID >= sequenceID {
			sequenceID += 1
		}
	}

	fmt.Println()
	//为请求分配编号
	request.SequenceID = sequenceID
	// 将请求消息存储到日志中
	state.MsgLogs.ReqMsg = request

	// hash  根据请求消息 生成hash
	digest, err := digest(request)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// 修改当前状态

	state.CurrentStage = PrePrepared
	// 返回预分配消息
	return &PrePrepareMsg{
		ViewID:     state.ViewID,
		SequenceID: sequenceID,
		Digest:     digest,
		RequestMsg: request,
	}, nil

}

// 序号分配
func (state *State) PrePrepare(prePrepareMsg *PrePrepareMsg) (*VoteMsg, error) {
	state.MsgLogs.ReqMsg = prePrepareMsg.RequestMsg

	if state.verifyMsg(prePrepareMsg.ViewID, prePrepareMsg.SequenceID, prePrepareMsg.Digest) {
		return nil, errors.New("pre-prepare message is corrupted ")
	}

	state.CurrentStage = PrePrepared
	return &VoteMsg{
		ViewID:     state.ViewID,
		SequenceID: prePrepareMsg.SequenceID,
		Digest:     prePrepareMsg.Digest,
		MsgType:    PrepareMsg,
	}, nil

}

// 相互交互
func (state *State) Prepare(prepareMsg *VoteMsg) (*VoteMsg, error) {

	if !state.verifyMsg(prepareMsg.ViewID, prepareMsg.SequenceID, prepareMsg.Digest) {
		return nil, errors.New("prepare message is corrupted")
	}
	state.MsgLogs.PrepareMsgs[prepareMsg.NodeID] = prepareMsg

	fmt.Printf("[Prepare-Vote]:%d \n", len(state.MsgLogs.PrepareMsgs))

	if state.prepared() {
		state.CurrentStage = Prepared

		return &VoteMsg{
			ViewID:     state.ViewID,
			SequenceID: prepareMsg.SequenceID,
			Digest:     prepareMsg.Digest,
			MsgType:    CommitMsg,
		}, nil
	}

	return nil, nil

}

// 提交节点
func (state *State) Commit(CommitMsg *VoteMsg) (*ReplyMsg, *RequestMsg, error) {

	if !state.verifyMsg(CommitMsg.ViewID, CommitMsg.SequenceID, CommitMsg.Digest) {
		return nil, nil, errors.New("prepare message is corrupted")
	}

	state.MsgLogs.CommitMsgs[CommitMsg.NodeID] = CommitMsg

	fmt.Printf("[Commit-Vote ]:%d \n", len(state.MsgLogs.CommitMsgs))

	if state.committed() {
		// This node executes the requested operation locally and gets the result.
		result := "Executed"

		// Change the stage to prepared.
		state.CurrentStage = Committed

		return &ReplyMsg{
			ViewID:    state.ViewID,
			Timestamp: state.MsgLogs.ReqMsg.Timestamp,
			ClientID:  state.MsgLogs.ReqMsg.ClientID,
			Result:    result,
		}, state.MsgLogs.ReqMsg, nil
	}

	return nil, nil, nil

}

func digest(object interface{}) (string, error) {
	msg, err := json.Marshal(object)
	if err != nil {
		return "", err
	}
	return Hash(msg), nil
}

func (state *State) verifyMsg(viewID int64, sequenceID int64, digestGot string) bool {
	// Wrong view. That is, wrong configurations of peers to start the consensus.
	if state.ViewID != viewID {
		return false
	}

	// Check if the Primary sent fault sequence number. => Faulty primary.
	// TODO: adopt upper/lower bound check.
	if state.LastSequenceID != -1 {
		if state.LastSequenceID >= sequenceID {
			return false
		}
	}

	digest, err := digest(state.MsgLogs.ReqMsg)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Check digest.
	if digestGot != digest {
		return false
	}

	return true
}

// 提交
func (state *State) committed() bool {
	if !state.prepared() {
		return false
	}
	if len(state.MsgLogs.CommitMsgs) < 2*f {
		return false
	}
	return true
}

func (state *State) prepared() bool {

	if state.MsgLogs.ReqMsg == nil {
		return false
	}

	if len(state.MsgLogs.PrepareMsgs) < 2*f {
		return false
	}

	return true

}
