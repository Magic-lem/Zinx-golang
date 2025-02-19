package ziface

// 消息管理抽象接口
type IMsgHandle interface {
	// 根据msgID索引执行路由方法，以非阻塞方式处理消息
	DoMsgHandler(request IRequest)
	// 添加路由方法到Map中
	AddRouter(msgId uint32, router IRouter)
	// 启动业务Worker工作池
	StartWorkerPool()
	// 将消息提交给TaskQueue的API
	SendMsgToTaskQueue(request IRequest)
}