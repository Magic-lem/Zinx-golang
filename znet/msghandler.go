package znet

import (
	"fmt"
	"strconv"
	"workspace/src/zinx/ziface"
)

type MsgHandle struct {
	// 保存 MsgId: router 的 Map
	Apis map[uint32] ziface.IRouter
}

/* --------- 实现各个方法 ------------ */
// 创建MsgHandle的方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32] ziface.IRouter),
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