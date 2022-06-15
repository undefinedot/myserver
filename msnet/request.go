package msnet

import "myserver/msiface"

type Request struct {
	conn msiface.IConnection // 已经于客户端建立的连接
	msg  msiface.IMessage    // 来自客户端的请求数据
}

func (r *Request) GetConnection() msiface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}
