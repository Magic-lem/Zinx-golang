package znet

import (
	"fmt"
	"strconv"
	"workspace/src/zinx/ziface"
	"workspace/src/zinx/utils"
)

type MsgHandle struct {
	// 保存 MsgId: router 的 Map
	Apis map[uint32] ziface.IRouter
	// 业务worker工作池的数量（协程数量）
	WorkerPoolSize uint32
	// 任务队列，用管道来实现
	TaskQueue []chan ziface.IRequest
}

/* --------- 实现各个方法 ------------ */
// 创建MsgHandle的方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32] ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue: make([]chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen),
	}
}

// 根据msgID索引执行路由方法，以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 1. 从request中找到msgId所对应的router
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgId(), " is Not found! Need Register!")
	}

	// 2. 调度router业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 添加路由方法到Map中
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	// 1. 判断当前msg绑定的API处理方法是否存在
	if _, ok := mh.Apis[msgId]; ok == true {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	// 2. 添加msg与API的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add Router msgId = ", msgId)
}

// 启动一个Worker，从TaskQueue中取出任务执行
func (mg *MsgHandle) StartOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerId, " is started.")
	// 无限循环，执行完一个任务再执行下一个，直到结束
	for {
		select {
		case request := <-taskQueue:
			// 有消息则取出队列的Request，并执行绑定的业务方法
			mg.DoMsgHandler(request)
		}
	}
}

// 启动业务Worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	// 遍历需要启动worker的数量，依此启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 初始化当前worker所对应的任务队列，开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)  
		// 利用Goroutine启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 将消息提交给TaskQueue的API
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 根据ConnID来分配当前的连接应该由哪个worker负责处理：轮询的平均分配法则
	// 得到需要处理此条连接的workerID
	workerId := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID()," request msgID=", request.GetMsgId(), "to workerID=", workerId)
	// 将消息发送给对应的消息队列
	mh.TaskQueue[workerId] <- request
}