package msnet

import "myserver/msiface"

type Request struct {
	conn msiface.IConnection // 已经于客户端建立的连接
	data []byte // 来自客户端的请求数据
}

func (r *Request) GetConnection() msiface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}
