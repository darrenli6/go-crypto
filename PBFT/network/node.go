package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"goland/go-crypto/PBFT/consensus"
)

type Node struct {
	// 节点名称
	NodeID string
	//key=nodeid value=url
	NodeTable map[string]string
	//节点视图
	View *View

	// 当前状态
	CurrentState *consensus.State

	//提交信息
	CommittedMsgs []*consensus.RequestMsg

	MsgBuffer *MsgBuffer

	//用于实现服务器控制终端消息流的传递
	//定义消息人口
	MsgEntrance chan interface{}
	//定义消息传送
	MsgDelivery chan interface{}

	// 警报
	Alarm chan bool
}

type MsgBuffer struct {
	ReqMsgs        []*consensus.RequestMsg
	PrePrepareMsgs []*consensus.PrePrepareMsg
	PrepareMsgs    []*consensus.VoteMsg
	CommitMsgs     []*consensus.VoteMsg
}

//视图
type View struct {
	ID      int64
	Primary string
}

const ResolvingTimeDuration = time.Millisecond * 1000

func NewNode(nodeID string) *Node {

	const viewID = 10000000000

	node := &Node{
		NodeID: nodeID,
		// 用于测试终端
		NodeTable: map[string]string{
			"N1": "localhost:5000",
			"N2": "localhost:5001",
			"N3": "localhost:5002",
			"N4": "localhost:5003",
		},
		View: &View{
			ID:      viewID,
			Primary: "N1",
		},
		CurrentState:  nil,
		CommittedMsgs: make([]*consensus.RequestMsg, 0),
		MsgBuffer: &MsgBuffer{
			ReqMsgs:        make([]*consensus.RequestMsg, 0),
			PrePrepareMsgs: make([]*consensus.PrePrepareMsg, 0),
			PrepareMsgs:    make([]*consensus.VoteMsg, 0),
			CommitMsgs:     make([]*consensus.VoteMsg, 0),
		},
		//chanels
		MsgEntrance: make(chan interface{}),
		MsgDelivery: make(chan interface{}),
		Alarm:       make(chan bool),
	}

	// 开启信息调度
	go node.dispatchMsg()

	// 开启警报触发器
	go node.alarmToDispatcher()

	// Start message resolver
	go node.resolveMsg()

	return node

}

//用于发送请求的方法
func (node *Node) Broadcast(msg interface{}, path string) map[string]error {
	errorMap := make(map[string]error)

	for nodeID, url := range node.NodeTable {
		if nodeID == node.NodeID {
			continue
		}

		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			errorMap[nodeID] = err
			continue
		}
		// localhost:5000
		// localhost:5001
		//localhost:5002

		send(url+path, jsonMsg)
	}

	if len(errorMap) == 0 {
		return nil
	} else {
		return errorMap
	}
}

func (node *Node) Reply(msg *consensus.ReplyMsg) error {
	// Print all committed messages.
	for _, value := range node.CommittedMsgs {
		fmt.Printf("Committed value: %s, %d, %s, %d", value.ClientID, value.Timestamp, value.Operation, value.SequenceID)
	}
	fmt.Print("\n")

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// Client
	send(node.NodeTable[node.View.Primary]+"/reply", jsonMsg)

	return nil
}

// GetReq can be called when the node's CurrentState is nil.
// Consensus start procedure for the Primary.
//
func (node *Node) GetReq(reqMsg *consensus.RequestMsg) error {
	LogMsg(reqMsg)

	// Create a new state for the new consensus.
	// 为当前节点创建新状态
	err := node.createStateForNewConsensus()
	if err != nil {
		return err
	}

	// Start the consensus process.
	prePrepareMsg, err := node.CurrentState.StartConsensus(reqMsg)

	if err != nil {
		return err
	}

	LogStage(fmt.Sprintf("Consensus Process (ViewID:%d)", node.CurrentState.ViewID), false)

	// Send getPrePrepare message
	if prePrepareMsg != nil {
		node.Broadcast(prePrepareMsg, "/preprepare")
		LogStage("Pre-prepare", true)
	}

	return nil
}

// GetPrePrepare can be called when the node's CurrentState is nil.
// Consensus start procedure for normal participants.
func (node *Node) GetPrePrepare(prePrepareMsg *consensus.PrePrepareMsg) error {
	LogMsg(prePrepareMsg)

	// Create a new state for the new consensus.
	err := node.createStateForNewConsensus()
	if err != nil {
		return err
	}

	prePareMsg, err := node.CurrentState.PrePrepare(prePrepareMsg)
	if err != nil {
		return err
	}

	if prePareMsg != nil {
		// Attach node ID to the message
		prePareMsg.NodeID = node.NodeID

		LogStage("Pre-prepare", true)
		node.Broadcast(prePareMsg, "/prepare")
		LogStage("Prepare", false)
	}

	return nil
}

func (node *Node) GetPrepare(prepareMsg *consensus.VoteMsg) error {
	LogMsg(prepareMsg)

	commitMsg, err := node.CurrentState.Prepare(prepareMsg)
	if err != nil {
		return err
	}

	if commitMsg != nil {
		// Attach node ID to the message
		commitMsg.NodeID = node.NodeID

		LogStage("Prepare", true)
		node.Broadcast(commitMsg, "/commit")
		LogStage("Commit", false)
	}

	return nil
}

func (node *Node) GetCommit(commitMsg *consensus.VoteMsg) error {
	LogMsg(commitMsg)

	replyMsg, committedMsg, err := node.CurrentState.Commit(commitMsg)
	if err != nil {
		return err
	}

	if replyMsg != nil {
		if committedMsg == nil {
			return errors.New("committed message is nil, even though the reply message is not nil")
		}

		// Attach node ID to the message
		replyMsg.NodeID = node.NodeID

		// Save the last version of committed messages to node.
		node.CommittedMsgs = append(node.CommittedMsgs, committedMsg)

		LogStage("Commit", true)
		node.Reply(replyMsg)
		LogStage("Reply", true)
	}

	return nil
}

func (node *Node) GetReply(msg *consensus.ReplyMsg) {
	fmt.Printf("Result: %s by %s\n", msg.Result, msg.NodeID)
}

func (node *Node) createStateForNewConsensus() error {
	// Check if there is an ongoing consensus process.
	// 检查现在是否正在进行共识
	if node.CurrentState != nil {
		return errors.New("another consensus is ongoing")
	}

	// Get the last sequence ID
	var lastSequenceID int64
	// 判断当前是否处理过请求消息
	if len(node.CommittedMsgs) == 0 {
		lastSequenceID = -1
	} else {
		// 取出本系统最后一次处理的消息
		lastSequenceID = node.CommittedMsgs[len(node.CommittedMsgs)-1].SequenceID
	}

	// Create a new state for this new consensus process in the Primary
	// 创建状态
	node.CurrentState = consensus.CreateState(node.View.ID, lastSequenceID)

	LogStage("Create the replica status", true)

	return nil
}

func (node *Node) dispatchMsg() {
	for {
		// 从MsgEntrance取消息
		select {

		case msg := <-node.MsgEntrance:
			err := node.routeMsg(msg)
			if err != nil {
				fmt.Println(err)
				// TODO: send err to ErrorChannel
			}

		case <-node.Alarm:
			// 接受警报
			err := node.routeMsgWhenAlarmed()
			if err != nil {
				fmt.Println(err)
				// TODO: send err to ErrorChannel
			}
		}
	}
}

func (node *Node) routeMsg(msg interface{}) []error {
	switch msg.(type) {
	// 请求消息
	case *consensus.RequestMsg:
		// 判断是否正在共识
		if node.CurrentState == nil {
			// Copy buffered messages first.
			// 创建切片 用于处理请求消息
			msgs := make([]*consensus.RequestMsg, len(node.MsgBuffer.ReqMsgs))

			// 将缓冲区的消息拷贝到切片中
			copy(msgs, node.MsgBuffer.ReqMsgs)

			// Append a newly arrived message.
			// 将最新的请求消息追加到切片
			msgs = append(msgs, msg.(*consensus.RequestMsg))

			// Empty the buffer.
			// 清空缓冲区的请求消息
			node.MsgBuffer.ReqMsgs = make([]*consensus.RequestMsg, 0)

			// Send messages.
			// 将消息发送到传送通道
			node.MsgDelivery <- msgs
		} else {
			// 正在共识，所以该消息无法处理 暂时放到将该消息存储缓冲区内
			node.MsgBuffer.ReqMsgs = append(node.MsgBuffer.ReqMsgs, msg.(*consensus.RequestMsg))
		}
	case *consensus.PrePrepareMsg:
		if node.CurrentState == nil {
			// Copy buffered messages first.
			msgs := make([]*consensus.PrePrepareMsg, len(node.MsgBuffer.PrePrepareMsgs))
			copy(msgs, node.MsgBuffer.PrePrepareMsgs)

			// Append a newly arrived message.
			msgs = append(msgs, msg.(*consensus.PrePrepareMsg))

			// Empty the buffer.
			node.MsgBuffer.PrePrepareMsgs = make([]*consensus.PrePrepareMsg, 0)

			// Send messages.
			// 发送到delivery
			node.MsgDelivery <- msgs
		} else {
			node.MsgBuffer.PrePrepareMsgs = append(node.MsgBuffer.PrePrepareMsgs, msg.(*consensus.PrePrepareMsg))
		}
	case *consensus.VoteMsg:
		// 区分一下 PrepareMsg CommitMsg
		if msg.(*consensus.VoteMsg).MsgType == consensus.PrepareMsg {

			if node.CurrentState == nil || node.CurrentState.CurrentStage != consensus.PrePrepared {
				node.MsgBuffer.PrepareMsgs = append(node.MsgBuffer.PrepareMsgs, msg.(*consensus.VoteMsg))
			} else {
				// Copy buffered messages first.
				msgs := make([]*consensus.VoteMsg, len(node.MsgBuffer.PrepareMsgs))
				copy(msgs, node.MsgBuffer.PrepareMsgs)

				// Append a newly arrived message.
				msgs = append(msgs, msg.(*consensus.VoteMsg))

				// Empty the buffer.
				node.MsgBuffer.PrepareMsgs = make([]*consensus.VoteMsg, 0)

				// Send messages.
				node.MsgDelivery <- msgs
			}
		} else if msg.(*consensus.VoteMsg).MsgType == consensus.CommitMsg {

			// 节点当前的状态为nil或者不是序号分配阶段，将相互信息存储到切片中
			// 因为节点当前的状态为nil,说当前没有进行共识 ，如果不是序号分配阶段 说明共识还没有进行到相互交互阶段

			if node.CurrentState == nil || node.CurrentState.CurrentStage != consensus.Prepared {
				node.MsgBuffer.CommitMsgs = append(node.MsgBuffer.CommitMsgs, msg.(*consensus.VoteMsg))
			} else {

				// Copy buffered messages first.
				// 创建切片

				msgs := make([]*consensus.VoteMsg, len(node.MsgBuffer.CommitMsgs))
				// 将缓冲区的消息复制到切片中中
				copy(msgs, node.MsgBuffer.CommitMsgs)

				// Append a newly arrived message.
				// 将最新的消息也追加到切片中
				msgs = append(msgs, msg.(*consensus.VoteMsg))

				// Empty the buffer.
				// 清空缓冲
				node.MsgBuffer.CommitMsgs = make([]*consensus.VoteMsg, 0)

				// Send messages.
				// 将消息发送到传送通道
				node.MsgDelivery <- msgs
			}
		}
	}

	return nil
}

func (node *Node) routeMsgWhenAlarmed() []error {
	if node.CurrentState == nil {
		// Check ReqMsgs, send them.
		if len(node.MsgBuffer.ReqMsgs) != 0 {
			msgs := make([]*consensus.RequestMsg, len(node.MsgBuffer.ReqMsgs))
			copy(msgs, node.MsgBuffer.ReqMsgs)

			node.MsgDelivery <- msgs
		}

		// Check PrePrepareMsgs, send them.
		if len(node.MsgBuffer.PrePrepareMsgs) != 0 {
			msgs := make([]*consensus.PrePrepareMsg, len(node.MsgBuffer.PrePrepareMsgs))
			copy(msgs, node.MsgBuffer.PrePrepareMsgs)

			node.MsgDelivery <- msgs
		}
	} else {
		switch node.CurrentState.CurrentStage {
		case consensus.PrePrepared:
			// Check PrepareMsgs, send them.
			if len(node.MsgBuffer.PrepareMsgs) != 0 {
				msgs := make([]*consensus.VoteMsg, len(node.MsgBuffer.PrepareMsgs))
				copy(msgs, node.MsgBuffer.PrepareMsgs)

				node.MsgDelivery <- msgs
			}
		case consensus.Prepared:
			// Check CommitMsgs, send them.
			if len(node.MsgBuffer.CommitMsgs) != 0 {
				msgs := make([]*consensus.VoteMsg, len(node.MsgBuffer.CommitMsgs))
				copy(msgs, node.MsgBuffer.CommitMsgs)

				node.MsgDelivery <- msgs
			}
		}
	}

	return nil
}

func (node *Node) resolveMsg() {
	for {
		// Get buffered messages from the dispatcher.
		// 从消息传送通道获取消息
		msgs := <-node.MsgDelivery
		switch msgs.(type) {

		case []*consensus.RequestMsg:
			errs := node.resolveRequestMsg(msgs.([]*consensus.RequestMsg))
			if len(errs) != 0 {
				for _, err := range errs {
					fmt.Println(err)
				}
				// TODO: send err to ErrorChannel
			}
		case []*consensus.PrePrepareMsg:
			errs := node.resolvePrePrepareMsg(msgs.([]*consensus.PrePrepareMsg))
			if len(errs) != 0 {
				for _, err := range errs {
					fmt.Println(err)
				}
				// TODO: send err to ErrorChannel
			}
		case []*consensus.VoteMsg:
			voteMsgs := msgs.([]*consensus.VoteMsg)
			if len(voteMsgs) == 0 {
				break
			}

			if voteMsgs[0].MsgType == consensus.PrepareMsg {
				errs := node.resolvePrepareMsg(voteMsgs)
				if len(errs) != 0 {
					for _, err := range errs {
						fmt.Println(err)
					}
					// TODO: send err to ErrorChannel
				}
			} else if voteMsgs[0].MsgType == consensus.CommitMsg {
				errs := node.resolveCommitMsg(voteMsgs)
				if len(errs) != 0 {
					for _, err := range errs {
						fmt.Println(err)
					}
					// TODO: send err to ErrorChannel
				}
			}
		}
	}
}

func (node *Node) alarmToDispatcher() {
	for {
		time.Sleep(ResolvingTimeDuration)
		node.Alarm <- true
	}
}

// 处理客户端的请求消息
func (node *Node) resolveRequestMsg(msgs []*consensus.RequestMsg) []error {
	errs := make([]error, 0)

	// Resolve messages
	// 遍历消息切片，取出消息
	for _, reqMsg := range msgs {
		err := node.GetReq(reqMsg)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}

func (node *Node) resolvePrePrepareMsg(msgs []*consensus.PrePrepareMsg) []error {
	errs := make([]error, 0)

	// Resolve messages
	for _, prePrepareMsg := range msgs {
		err := node.GetPrePrepare(prePrepareMsg)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}

func (node *Node) resolvePrepareMsg(msgs []*consensus.VoteMsg) []error {
	errs := make([]error, 0)

	// Resolve messages
	for _, prepareMsg := range msgs {
		err := node.GetPrepare(prepareMsg)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}

func (node *Node) resolveCommitMsg(msgs []*consensus.VoteMsg) []error {
	errs := make([]error, 0)

	// Resolve messages
	for _, commitMsg := range msgs {
		err := node.GetCommit(commitMsg)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}
