package ziface

/*
	封包数据和拆包数据
	直接面向TCP连接中的数据流,为传输数据添加头部信息，用于处理TCP粘包问题。
*/
type IDataPack interface {
	GetHaedLen() uint32   // 获取包头部的长度
	Pack(IMessage) ([]byte, error)  // 将输入的消息封包
	UnPack([]byte) (IMessage, error)	   // 拆包			
}