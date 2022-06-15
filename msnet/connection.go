package msnet

import (
	"fmt"
	"myserver/msiface"
	"myserver/utils"
	"net"
)

type Connection struct {
	Conn *net.TCPConn // 连接的socket
	ConnID uint32 // 连接的ID
	isClosed bool // 连接状态
	ExitChan chan bool // 告知连接已经退出，缓冲区为1的channel TODO: 使用context

	Router msiface.IRouter // 处理连接的业务的方法从路由获取
}

func NewConnection(conn *net.TCPConn, connID uint32, router msiface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		Router: router,
	}
	return c
}
// StartReader 开启读数据业务
func (c *Connection) StartReader() {
	fmt.Println("Reader goroutine is running...")

	defer fmt.Println("log: connID =", c.ConnID, "Reader is exit, remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// read data
		buf := make([]byte, utils.GlobalConfig.MaxPacketSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("Read buf err:", err)
			continue
		}

		// TODO: req是否必须是值，再传参时取地址？
		// 封装Request，Router处理业务，将Request传入Router的方法中
		req := &Request{
			conn: c,
			data: buf,
		}
		go func(r msiface.IRequest) {
			// todo: [记]模板方法设计模式，固定执行顺序，用户只能改写方法的具体实现
			c.Router.PreHandle(r)
			c.Router.Handle(r)
			c.Router.PostHandle(r)
		}(req)
	}
}

// StartWriter 开启写数据业务
func (c *Connection) StartWriter() {

}
func (c *Connection) Start() {
	fmt.Println("Conn Start, ConnID is", c.ConnID)

	// 分别开启goroutine来read和write
	go c.StartReader()
	go c.StartWriter()
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop, ConnID is", c.ConnID)

	// TODO：什么时候在调用stop前已经是true
	if c.isClosed {
		return
	}
	c.isClosed = true

	// close socket
	c.Conn.Close()
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	return nil
}
