package hiface

type Iserver interface {
	//启动服务器
	Start()

	//停止服务器
	Stop()

	//运行服务器
	Serve()

	AddRouter(msgID uint32, router IRouter)
}
