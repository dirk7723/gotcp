package hiface

type IRequest interface {
	GetConnection() IConnection

	GetData() []byte

	GetMsgId() uint32
}
