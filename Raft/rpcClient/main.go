package main

import (
	"fmt"
	"log"
	"net/rpc"
)

// 实现rpc客户端

type Params struct {
	Width, Height int
}

// 调用服务
func main() {
	// 链接远程rpc服务
	rpc, err := rpc.DialHTTP("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}

	// 定义接受服务端传递的结果
	ret := 0
	// 调用远程求面积方法

	errr2 := rpc.Call("Rect.Area", Params{20, 40}, &ret)

	if errr2 != nil {
		log.Fatal(errr2)
	}

	fmt.Println("面积：", ret)
	err2 := rpc.Call("Rect.Perimeter", Params{10, 40}, &ret)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("周长：", ret)

}
