package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"goland/go-crypto/PBFT/consensus"
)

type Server struct {
	url  string
	node *Node
}

func NewServer(nodeID string) *Server {
	//创建终端节点的对象
	node := NewNode(nodeID)
	server := &Server{node.NodeTable[nodeID], node}
	//触发回调函数
	server.setRoute()
	return server
}

func (server *Server) Start() {
	fmt.Printf("服务启动在 %s...\n", server.url)
	//根据url做监听
	if err := http.ListenAndServe(server.url, nil); err != nil {
		//打印错误信息
		fmt.Println(err)
		return
	}
}

//回调函数
func (server *Server) setRoute() {
	//发送的请求是http://localhost:1111/req，这里接收到
	http.HandleFunc("/req", server.getReq)
	http.HandleFunc("/preprepare", server.getPrePrepare)
	http.HandleFunc("/prepare", server.getPrepare)
	http.HandleFunc("/commit", server.getCommit)
	http.HandleFunc("/reply", server.getReply)
}

//获得请求数据
func (server *Server) getReq(writer http.ResponseWriter, request *http.Request) {
	//创建请求消息对象
	var msg consensus.RequestMsg

	//对数据解码
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	//将解析出的数据放到server端的通道中
	//为Server中的node赋值
	server.node.MsgEntrance <- &msg
}

func (server *Server) getPrePrepare(writer http.ResponseWriter, request *http.Request) {
	var msg consensus.PrePrepareMsg

	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	server.node.MsgEntrance <- &msg
}

func (server *Server) getPrepare(writer http.ResponseWriter, request *http.Request) {
	var msg consensus.VoteMsg
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	server.node.MsgEntrance <- &msg
}

func (server *Server) getCommit(writer http.ResponseWriter, request *http.Request) {
	var msg consensus.VoteMsg
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	server.node.MsgEntrance <- &msg
}

func (server *Server) getReply(writer http.ResponseWriter, request *http.Request) {
	var msg consensus.ReplyMsg
	err := json.NewDecoder(request.Body).Decode(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	server.node.GetReply(&msg)
}

func send(url string, msg []byte) {
	buff := bytes.NewBuffer(msg)

	http.Post("http://"+url, "application/json", buff)
}
