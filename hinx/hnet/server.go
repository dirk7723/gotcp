package hnet

import (
	"fmt"
	"net"

	"hbq.com/ggame/hinx/hiface"
	"hbq.com/ggame/hinx/untils"
)

type Server struct {
	//服务器名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的ip
	IP string
	//服务器监听的端口
	Port int

	//Router hiface.IRouter

	//消息的管理msgid和对应的处理业务api之间的关系
	Msghandle hiface.ImsgHandle
}

/* func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//fmt.Println("[Conn Handle] CallBackToClient")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("wirte buf error ", err)
		return errors.New("CallBackToClient error")
	}
	return nil
} */

func (s *Server) Start() {
	fmt.Printf("[start] server listenner at ip :%s port :%d, is starting\n", s.IP, s.Port)

	go func() {
		//开启消息队列及worker工作池
		s.Msghandle.StartWorkerPool()

		//1.获取一个TCP的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error : ", err)
			return
		}

		//2.监听服务器地址
		listen, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " error ", err)
			return
		}
		fmt.Println("start Hinx server succ, ", s.Name, " succ, listening...")

		/*----------------------------*/
		var cid uint32
		cid = 0
		//3.阻塞客户端连接，处理客户端连接业务（读写）
		for {
			conn, err := listen.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error ", err)
				continue
			}

			dealConn := NewConnection(conn, cid, s.Msghandle)
			cid++

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	//TODO
}

func (s *Server) Serve() {
	s.Start()

	//TODO 其他一些任务

	//阻塞状态
	select {}
}

func NewServer(name string) hiface.Iserver {

	s := &Server{
		Name:      untils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        untils.GlobalObject.Host,
		Port:      untils.GlobalObject.TcpPort,
		Msghandle: NewMsgHandle(),
	}

	return s
}

func (s *Server) AddRouter(msgID uint32, router hiface.IRouter) {
	s.Msghandle.AddRouter(msgID, router)
	fmt.Println("add router succ")
}
