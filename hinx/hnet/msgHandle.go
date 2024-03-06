package hnet

import (
	"fmt"
	"strconv"

	"hbq.com/ggame/hinx/hiface"
	"hbq.com/ggame/hinx/untils"
)

type Msghandle struct {
	Apis map[uint32]hiface.IRouter

	//负责worker取任务的消息队列
	TaskQueue []chan hiface.IRequest

	//业务工作worker池的数量
	WorkerPoolSize uint32
}

func NewMsgHandle() *Msghandle {
	return &Msghandle{
		Apis:           make(map[uint32]hiface.IRouter),
		WorkerPoolSize: untils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan hiface.IRequest, untils.GlobalObject.WorkerPoolSize),
	}
}

func (mh *Msghandle) DoMsgHandle(request hiface.IRequest) {
	handle, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgid = ", request.GetMsgId(), " is not found need register")
	}
	handle.PreHandle(request)
	handle.Handle(request)
	handle.PostHandle(request)
}

func (mh *Msghandle) AddRouter(msgID uint32, router hiface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api,msgid = " + strconv.Itoa(int(msgID)))
	}
	mh.Apis[msgID] = router
	fmt.Println("add api msgID = ", msgID, " success")
}

// 启动一个worker工作池
func (mh *Msghandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//1.当前的worker对应的channel消息队列 开辟空间 第0个worker 就用第0个channel
		mh.TaskQueue[i] = make(chan hiface.IRequest, untils.GlobalObject.MaxWorkerTaskLen)

		//启动当前的worker 阻塞等待消息从channel传递进来
		go mh.StartOneWoker(i, mh.TaskQueue[i])

	}
}

// 启动一个worker工作流程
func (mh *Msghandle) StartOneWoker(workID int, taskQueue chan hiface.IRequest) {
	fmt.Println("work id = ", workID, " is started...")
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandle(request)
		}
	}
}

// 将消息发送给消息任务队列
func (mh *Msghandle) SendMsgToTaskQueue(request hiface.IRequest) {
	//将消息平均分配给不同的worker
	workerID := request.GetConnection().GetTcpConnID() % mh.WorkerPoolSize
	fmt.Println("Add connID = ", request.GetConnection().GetTcpConnID(), " request msgid = ", request.GetMsgId(), " to workerID = ", workerID)

	//发给指定workerid的taskQueue
	mh.TaskQueue[workerID] <- request
}
