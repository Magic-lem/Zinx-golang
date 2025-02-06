package znet

import (
	"workspace/src/zinx/ziface"
)

type Request struct {
	conn ziface.IConnection   // 已经和客户端建立好的连接对象
	data []byte // 来自客户端的请求数据
}


// 获取当前连接
func (r *Request) GetConnetcion() ziface.IConnection {
	return r.conn
}

// 获取当前请求消息数据
func (r *Request) GetData() []byte {
	return r.data
}