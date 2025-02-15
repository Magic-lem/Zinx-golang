package ziface

type IMessage interface {
	GetDataLen() uint32  // 获取消息的长度
	GetMsgId()   uint32  // 获取消息的ID
	GetData() 	 []byte  // 获取消息

	SetMsgId(uint32)     // 设置消息的ID
	SetData([]byte)		 // 设置消息
	SetDataLen(uint32)   // 设置消息的长度
}