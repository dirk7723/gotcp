package main

import (
	"fmt"

	"hbq.com/ggame/hinx/hiface"
	"hbq.com/ggame/hinx/hnet"
)

type PingRoutrer struct {
	hnet.BaseRouter
}

func (pr *PingRoutrer) Handle(request hiface.IRequest) {
	fmt.Println("call back ping-Handle...")
	/* _, err := request.GetConnection().GetTcpConnection().Write([]byte("ping...\n"))
	if err != nil {
		fmt.Println("call back ping... error")
	} */

	fmt.Println("recv from client: msgid = ", request.GetMsgId(), ",data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloRoutrer struct {
	hnet.BaseRouter
}

func (pr *HelloRoutrer) Handle(request hiface.IRequest) {
	fmt.Println("call back hello-Handle...")
	/* _, err := request.GetConnection().GetTcpConnection().Write([]byte("ping...\n"))
	if err != nil {
		fmt.Println("call back ping... error")
	} */

	fmt.Println("recv from client: msgid = ", request.GetMsgId(), ",data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("hello...hello...hello..."))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := hnet.NewServer("[hinx V0.8]")

	s.AddRouter(0, &PingRoutrer{})

	s.AddRouter(1, &HelloRoutrer{})

	s.Serve()
}
