package msiface

// IMessage 封装请求的消息
type IMessage interface {
	GetMsgID() uint32
	GetDataLen() uint32
	GetData() []byte

	SetMsgID(uint32)
	SetDataLen(uint32)
	SetData([]byte)
}
