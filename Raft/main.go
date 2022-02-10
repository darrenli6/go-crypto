package main

import "sync"

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
	rf := Raft{}
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

	rf.electCh = make(chan  bool)


}

func (rf *Raft) setTerm(term int) {

	rf.currentTerm = term
}
