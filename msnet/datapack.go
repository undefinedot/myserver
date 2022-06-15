package msnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"myserver/msiface"
	"myserver/utils"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	// 固定长度
	return 8
}

// Pack 封包: Message=>[]byte类型的data
func (dp *DataPack) Pack(msg msiface.IMessage) ([]byte, error) {
	dataBuf := bytes.NewBuffer([]byte{})
	// 按顺序写，len
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	// ID
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	// data
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuf.Bytes(), nil
}

// Unpack 拆包: []byte=>Message; 只需读消息头,data为nil,从conn按消息头读数据
func (dp *DataPack) Unpack(data []byte) (msiface.IMessage, error) {
	reader := bytes.NewReader(data)
	msg := &Message{}

	// 按顺序读，len (参数3必须是指针)
	if err := binary.Read(reader, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// ID
	if err := binary.Read(reader, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}
	// data，需要限制大小
	maxSize := utils.GlobalConfig.MaxPacketSize
	if (maxSize > 0) && msg.DataLen > maxSize {
		return nil, errors.New("数据过大")
	}

	return msg, nil
}
