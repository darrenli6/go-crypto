package main

import (
	"log"
	"net/http"
	"net/rpc"
)

// 声明函数 必须符合RPC条件和标准
// 函数不能随便定义
// A和B传参数，参数类型必须一致
// go对rpc的支持，支持三个级别 TCP JSON HTTP
// go的RPC支持开发的服务端与客户端之间的交互

// 不能随便定义
// 首字母必须大写
type Params struct {
	// 参数必须首字母大写，要进行跨域访问的
	Width, Height int
}

// 声明一个矩形

type Rect struct {
}

// 函数必须导出的
// 必须有两个导出类型的参数
// 第一个参数是接受参数 首字母是大写
// 第二个参数返回客户端的参数，必须是指针类型
// 函数必须有一个返回值 error

// 求矩形面积
func (r *Rect) Area(p Params, ret *int) error {

	*ret = p.Width * p.Height
	return nil
}

// 求周长
func (r *Rect) Perimeter(p Params, ret *int) error {

	*ret = (p.Width + p.Height) * 2
	return nil
}

func main() {

	// 注册服务
	rect := new(Rect)
	// 注册一个rect服务
	// 封装的socket
	rpc.Register(rect)
	// 将服务处理绑定到http协议上
	rpc.HandleHTTP()
	// 监听
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
