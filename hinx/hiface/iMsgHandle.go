package hiface

type ImsgHandle interface {
	//执行对应的Router消息处理方法
	DoMsgHandle(request IRequest)

	//为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)

	StartWorkerPool()

	SendMsgToTaskQueue(request IRequest)
}
