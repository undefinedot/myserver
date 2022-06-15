package msnet

import (
	"errors"
	"fmt"
	"io"
	"myserver/msiface"
	"net"
)

type Connection struct {
	Conn     *net.TCPConn // 连接的socket
	ConnID   uint32       // 连接的ID
	isClosed bool         // 连接状态
	ExitChan chan bool    // 告知连接已经退出，缓冲区为1的channel TODO: 使用context

	MsgHandler msiface.IMsgHandler // 处理连接的业务的方法从路由获取
}

func NewConnection(conn *net.TCPConn, connID uint32, routers msiface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		MsgHandler: routers,
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

		// 用来处理包的对象
		dp := NewDataPack()
		// 读取 // TODO: headData每次都创建，有必要改为单例吗
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.Conn, headData); err != nil {
			fmt.Println("ReadFull headData error:", err)
			return //直接退出
		}
		// unpack
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("Unpack headData error:", err)
			return
		}

		// 根据head读data // TODO: msg.(*Message).Data是否改变接口变量?可能报错！
		if _, err := io.ReadFull(c.Conn, msg.(*Message).Data); err != nil {
			fmt.Println("ReadFull data error:", err)
			return
		}

		// 放入Request
		// TODO: req是否必须是值，再传参时取地址？
		// 封装Request，Router处理业务，将Request传入MsgHandler这个路由集合中
		req := &Request{
			conn: c,
			msg:  msg,
		}

		go c.MsgHandler.DoMsgHandler(req)
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

// SendMsg 先封包，再发数据
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	// 连接是否关闭
	if c.isClosed {
		return errors.New("Connection closed...")
	}
	// 封包
	dp := NewDataPack()
	msgByte, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error:", err)
		return err
	}
	if _, err := c.Conn.Write(msgByte); err != nil {
		fmt.Println("conn write error:", err)
		return err
	}

	return nil
}
