package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"hbq.com/ggame/hinx/hnet"
)

func main() {
	fmt.Println("client start...")

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start error ", err)
		return
	}

	for {
		/* fmt.Println("client conn ", conn)
		_, err := conn.Write([]byte("i am client ..."))
		if err != nil {
			fmt.Println("write conn error ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error ", err)
			return
		}

		fmt.Printf("server call back: %s,cnt: %d\n", buf, cnt) */

		dp := hnet.NewDataPack()

		binary_msg, err := dp.Pack(hnet.NewMsgPackage(0, []byte("Hinxv0.6 client test message")))
		if err != nil {
			fmt.Println("pack error:", err)
			return
		}
		if _, err := conn.Write(binary_msg); err != nil {
			fmt.Println("write error:", err)
			return
		}

		binary_head := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binary_head); err != nil {
			fmt.Println("read head error ", err)
			break
		}

		msg_head, err := dp.Unpack(binary_head)
		if err != nil {
			fmt.Println("client unpack msghead error", err)
			break
		}
		if msg_head.GetMsgLen() > 0 {
			msg := msg_head.(*hnet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error", err)
				return
			}

			fmt.Println("recv msgid:", msg.Id, ",datalen=", msg.DataLen, ",data=", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}

}
