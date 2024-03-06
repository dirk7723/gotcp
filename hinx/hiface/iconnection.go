package hiface

import "net"

type IConnection interface {
	Start()

	Stop()

	//获取当前连接绑定的socket conn
	GetTcpConnection() *net.TCPConn

	//获取当前连接模块的连接ID
	GetTcpConnID() uint32

	//获取远程客户端的 tcp状态 ip port
	RemoteAddr() net.Addr

	//发送数据
	SendMsg(msgId uint32, data []byte) error
}

type HandFunc func(*net.TCPConn, []byte, int) error
