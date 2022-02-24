package consensus

type RequestMsg struct {
	// 时间
	Timestamp int64 `json:"timestamp"`
	// 客户端编号
	ClientID string `json:"clientID"`
	// 客户端请求的具体信息
	Operation string `json:"operation"`
	//请求编号
	SequenceID int64 `json:"sequenceID"`
}

type ReplyMsg struct {
	ViewID int64 `json:"viewID"`
	//时间戳
	Timestamp int64 `json:"timstamp"`

	ClientID string `json:"clientID"`
	NodeID   string `json:"nodeID"`
	Result   string `json:"result"`
}

// 序号分配消息
type PrePrepareMsg struct {
	// 视图id
	ViewID int64 `json:"viewID"`

	// 请求id
	SequenceID int64 `json:"sequenceID"`
	//hash值
	Digest string `json:"digest"`

	// 请求信息
	RequestMsg *RequestMsg `json:"requestMsg"`
}

//
type VoteMsg struct {
	// 视图id
	ViewID int64 `json:"viewID"`
	// 请求ID
	SequenceID int64 `json:"sequenceID"`
	// hash值
	Digest string `json:"digest"`
	// 发送信息的编号
	NodeID string `json:"nodeID"`
	//
	MsgType `json:"msgType"`
}

type MsgType int

const (
	PrepareMsg MsgType = iota
	CommitMsg
)
