package hnet

import "hbq.com/ggame/hinx/hiface"

type Request struct {
	conn hiface.IConnection
	msg  hiface.IMessage
}

func (r *Request) GetConnection() hiface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}
