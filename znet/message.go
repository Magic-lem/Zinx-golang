package znet

import (
	"workspace/src/zinx/ziface"
)

type Message struct {
	Id        uint32     // 消息的ID
	DataLen   uint32	 // 消息的长度
	Data	  []byte     // 消息的内容
}
/* -------------- 构造函数 --------------------*/
func NewMsgPackage(id uint32, data []byte) ziface.IMessage {
	return &Message{
		Id: id,
		Data: data,
		DataLen: uint32(len(data)),
	}
}

/* -------------- 实现各个方法 ----------------- */
// 获取消息的长度
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}
// 获取消息的ID
func (m *Message) GetMsgId() uint32 {
	return m.Id
} 
// 获取消息
func (m *Message) GetData() []byte {
	return m.Data
}
// 设置消息的ID
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}
// 设置消息 
func (m *Message) SetData(data []byte) {
	m.Data = data
}
// 设置消息的长度 
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}