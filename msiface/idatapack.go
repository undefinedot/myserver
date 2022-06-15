package msiface

type IDataPack interface {
	GetHeadLen() uint32 // 获取消息头长度
	Pack(msg IMessage) ([]byte, error) // 封包
	Unpack([]byte) (IMessage, error) // 解包
}
