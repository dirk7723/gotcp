package hnet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"hbq.com/ggame/hinx/hiface"
	"hbq.com/ggame/hinx/untils"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	return 8 //datalen uint32(4字节) + ID uint32(4字节)
}

func (dp *DataPack) Pack(msg hiface.IMessage) ([]byte, error) {
	//创还能一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (dp *DataPack) Unpack(binary_data []byte) (hiface.IMessage, error) {
	//创建一个从输入二进制数据的ioreader
	bytes_reader := bytes.NewReader(binary_data)

	msg := &Message{}

	if err := binary.Read(bytes_reader, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(bytes_reader, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	if untils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > untils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv")
	}

	return msg, nil
}
