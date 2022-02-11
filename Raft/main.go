package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 模拟选举的逻辑
// 模拟三节点的分布式选举

// 定义常量

const raftCount = 3

type Leader struct {
	// 任期
	Term int
	// 领导编号
	LeaderId int
}

// 创建存储leader的对象
// 最初任期为0 -1代表没有编号

var leader = Leader{0, -1}

//声明raft节点类型
type Raft struct {
	// 锁
	mu sync.Mutex

	// 节点编号
	me int
	// 当前任期
	currentTerm int
	// 为哪个节点投票
	votedFor int
	// 当前节点的状态
	// 0 follower 1 candidate 2 leader
	state int
	// 发送最后一条消息的时间
	lastMessageTime int64
	// 当前节点的领导
	currentLeader int

	// 消息通道
	message chan bool

	// 选举通道
	electCh chan bool

	// 心跳信号
	heartBeat chan bool

	// 返回心跳信号
	hearbeatRe chan bool

	//超时时间
	timeout int
}

func main() {

	// 创建三个节点 最初是follower 状态
	// 如果出现candidate状态 开始投票
	// 最后产生leader

	// 创建三个节点
	for i := 0; i < raftCount; i++ {
		//创建节点的make
		Make(i)
	}
}

// 创建节点
func Make(me int) *Raft {
	rf := &Raft{}
	// 编号
	rf.me = me
	// 0 1 2 -1 都不投票
	rf.votedFor = -1

	//0 是folower状态
	rf.state = 0
	rf.timeout = 0

	// 最初没有领导
	rf.currentLeader = -1

	//设置任期
	rf.setTerm(0)

	// 通道
	rf.electCh = make(chan bool)
	rf.message = make(chan bool)
	rf.heartBeat = make(chan bool)

	rf.hearbeatRe = make(chan bool)

	//随机种子
	rand.Seed(time.Now().UnixNano())

	// 选举逻辑
	go rf.election()
	// 心跳检查
	go rf.sendLeaderHeartBeat()

	return rf

}

func (rf *Raft) election() {

	// 设置
	var result bool
	// 循环投票
	for {
		timeout := randRange(150, 900)

		// 设置最后一条消息的时间
		rf.lastMessageTime = millisecond()
		select {
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			fmt.Println("当前节点状态是", rf.state)
		}

		result = false
		// 选出leader，停止循环，result=true
		for !result {
			// 选择leader
			result = rf.election_ont_out(&leader)
		}
	}
}

func (rf *Raft) sendLeaderHeartBeat() {

}

func (rf *Raft) setTerm(term int) {

	rf.currentTerm = term

}

//生产随机数
func randRange(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}

// 获取当前时间
func millisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// 选举leader
func (rf *Raft) election_ont_out(leader Leader) bool {

	// 超时时间
	var timeout int64
	timeout = 100
	// 投票数量
	var vote int64
	// 是否开启心跳方法
	var triggerHeartbeat bool
	// 当前时间戳的
	last := millisecond()
	// 定义返回值

	success := false

	// 首先成为condidate状态
	rf.mu.Lock()
	rf.becomeCondidate()
	rf.mu.Unlock()

	// 开始选举
	fmt.Println("start election leader")

	for {
		// 便利所有节点进行投票
		for i := 0; i < raftCount; i++ {
			// 遍历到不是自己 进行来票
			if i != rf.me {
				go func() {
					if leader.LeaderId < 0 {
						// 其他节点没有领导
						rf.electCh <- true

					}
				}()
			}

		}

		//设置投票数量
		vote = 0
		triggerHeartbeat = false
		// 遍历所有节点进行选举
		for i := 0; i < raftCount; i++ {
			select {
			case ok := <-rf.electCh:
				if ok {

				}
			}
		}

	}
	return success

}

// 修改节点为candidate
func (rf *Raft) becomeCondidate() {
	// 设置状态
	rf.state = 1
	//节点的任期+1
	rf.setTerm(rf.currentTerm + 1)

	// 设置为哪个节点投票
	rf.votedFor = rf.me

	// 是否有领导
	rf.currentLeader = -1

}
