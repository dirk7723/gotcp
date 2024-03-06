package hnet

import (
	"errors"
	"fmt"
	"io"
	"net"

	"hbq.com/ggame/hinx/hiface"
	"hbq.com/ggame/hinx/untils"
)

type Connection struct {
	Conn *net.TCPConn

	ConnID uint32

	isClosed bool

	//handAPI hiface.HandFunc

	ExitChan chan bool

	//该连接处理的方法router
	//Router hiface.IRouter

	//消息的管理msgid和对应的处理业务api之间的关系
	Msghandle hiface.ImsgHandle

	//
	msgChan chan []byte
}

func NewConnection(conn *net.TCPConn, connID uint32, msghandle hiface.ImsgHandle) *Connection {
	c := &Connection{
		Conn:   conn,
		ConnID: connID,
		//Router:   router,
		isClosed:  false,
		ExitChan:  make(chan bool, 1),
		Msghandle: msghandle,
		msgChan:   make(chan []byte),
	}
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("[StartReader is running]")
	defer fmt.Println("connid ", c.ConnID, " reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		/* buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf error ", err)
			continue
		} */

		//创建一个拆包解包对象
		dp := NewDataPack()
		head_data := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTcpConnection(), head_data); err != nil {
			fmt.Println("read msg head error", err)
			break
		}

		//拆包，得到msgid和msgdatalen,放在msg消息中
		msg, err := dp.Unpack(head_data)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTcpConnection(), data); err != nil {
				fmt.Println("read msg data error:", err)
				break
			}
		}
		msg.SetData(data)

		//调用当前连接所绑定的HandleApi
		/* if err := c.handAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("connid ", c.ConnID, "hand is error", err)
			break
		}
		fmt.Println("Conn ", c.Conn)
		fmt.Printf("client send is %s\n", buf[:cnt]) */

		req := &Request{
			conn: c,
			msg:  msg,
		}

		/* c.Router.PreHandle(req)

		c.Router.Handle(req)

		c.Router.PostHandle(req) */

		//go c.Msghandle.DoMsgHandle(req)

		if untils.GlobalObject.MaxPackageSize > 0 {
			c.Msghandle.SendMsgToTaskQueue(req) //工作池处理
		} else {
			go c.Msghandle.DoMsgHandle(req) //原生gorouteine处理
		}
	}

}

// 写消息Gorouteine 专门发送给客户端消息的模块
func (c *Connection) StartWrite() {
	fmt.Println("[write gorouteine is running]")
	defer fmt.Println(c.Conn.RemoteAddr().String(), " [conn write exit]")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error: ", err)
				return
			}

		case <-c.ExitChan:
			//代表reader已经退出，writer也需要退出
			return
		}
	}
}

func (c *Connection) Start() {

	fmt.Println("Conn Start().. connid=", c.ConnID)
	//当前连接读
	go c.StartReader()
	//当前连接写
	go c.StartWrite()
}

func (c *Connection) Stop() {
	fmt.Println("Conn stop().. connid=", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true
	//关闭socket连接
	c.Conn.Close()

	//告知writer关闭
	c.ExitChan <- true

	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

// 获取当前连接绑定的socket conn
func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接模块的连接ID
func (c *Connection) GetTcpConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的 tcp状态 ip port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed msg")
	}
	//将data进行封包 msgdatalen msgid data
	dp := NewDataPack()
	binary_msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack error msg id=", msgId)
		return errors.New("pack erros msg")
	}

	//将数据发送给客户端
	/* if _, err := c.Conn.Write(binary_msg); err != nil {
		fmt.Println("write msg id=", msgId, "error:", err)
		return errors.New("pack erros msg")
	} */

	//将数据发送给管道
	c.msgChan <- binary_msg

	return nil
}
