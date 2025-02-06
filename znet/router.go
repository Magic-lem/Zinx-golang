package znet

import (
	"workspace/src/zinx/ziface"
)

// 定义一个基类，后面实现Router时需要继承这个基类，根据需求对基类的方法进行重写就好了
type BaseRouter struct {}

// 下面函数都没有实现，因为这是基类，后面具体的Router类需要重写
// 在处理conn业务之前的钩子方法Hook
func (br *BaseRouter) PreHandler(request ziface.IRequest) {}
// 处理coon业务的主方法Hook
func (br *BaseRouter) Handler(request ziface.IRequest) {}
// 在处理conn业务之后的钩子方法Hook
func (br *BaseRouter) PostHandler(request ziface.IRequest) {}

