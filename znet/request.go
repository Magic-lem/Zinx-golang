package znet

import (
	"workspace/src/zinx/ziface"
)

type Request struct {
	conn ziface.IConnection   // 已经和客户端建立好的连接对象
	msg ziface.IMessage		  // 客户端请求消息        update：Zinx-V0.5修改为封装好的消息类型
}


// 获取当前连接
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// 获取当前请求消息数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()   // update：Zinx-V0.5修改为从消息中提取
}

// Zinx-V0.5 获取当前请求小的ID
func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()   
}