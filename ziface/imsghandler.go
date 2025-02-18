package ziface

// 消息管理抽象接口
type IMsgHandle interface {
	// 根据msgID索引执行路由方法，以非阻塞方式处理消息
	DoMsgHandler(request IRequest)
	// 添加路由方法到Map中
	AddRouter(msgId uint32, router IRouter)
}