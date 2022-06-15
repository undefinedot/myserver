package msiface

// IRequest 封装Connection和请求数据在一个Request中
type IRequest interface {
	GetConnection() IConnection // TODO: 值or指针
	GetData() []byte
	GetMsgID() uint32
}
