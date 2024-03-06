package hiface

type Idatapack interface {
	GetHeadLen() uint32

	Pack(msg IMessage) ([]byte, error)

	unpack([]byte) (IMessage, error)
}
