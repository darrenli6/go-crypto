package consensus

type PBFT interface {

	// 开始共识
	StartConsensus(request *RequestMsg) (*PrePrepareMsg, error)

	// 序号分配
	PrePrepare(PrePrepareMsg *PrePrepareMsg) (*VoteMsg, error)

	//相互交互
	Prepare(prepareMsg *VoteMsg) (*VoteMsg, error)

	// 序号确认
	// 会对requestMsg校验 是否合法
	Commit(CommitMsg *VoteMsg) (*ReplyMsg, *RequestMsg, error)
}
