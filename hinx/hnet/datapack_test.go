package hnet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		return
	}

	go func() {
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accpt  error", err)
			}

			go func(conn net.Conn) {
				dp := NewDataPack()
				for {
					headDta := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headDta)
					if err != nil {
						fmt.Println("read head error")
						break
					}

					msg_head, err := dp.Unpack(headDta)
					if err != nil {
						fmt.Println("server unpack error")
						return
					}

					if msg_head.GetMsgLen() > 0 {
						//msg是有效数据，需要读取

						//
						msg := msg_head.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack error")
							return
						}
						fmt.Println("recv msgid:", msg.Id, ",datalen=", msg.DataLen, ",data=", msg.Data)

					}

				}
			}(conn)
		}
	}()

	//模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client Dial error:", err)
		return
	}

	dp := NewDataPack()
	//模拟2次粘包
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte("hinx"),
	}
	send_data1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error:", err)
		return
	}

	msg2 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte("hello"),
	}
	send_data2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error:", err)
		return
	}

	send_data1 = append(send_data1, send_data2...)
	conn.Write(send_data1)

	select {}
}
