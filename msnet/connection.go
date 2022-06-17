package msnet

import (
	"errors"
	"fmt"
	"io"
	"myserver/msiface"
	"myserver/utils"
	"net"
	"sync"
)

type Connection struct {
	TcpServer  msiface.IServer
	Conn       *net.TCPConn        // 连接的socket
	ConnID     uint32              // 连接的ID
	isClosed   bool                // 连接状态
	ExitChan   chan bool           // 告知连接已经退出，缓冲区为1的channel TODO: 使用context
	msgChan    chan []byte         // reader和writer两个goroutine通信使用
	MsgHandler msiface.IMsgHandler // 处理连接的业务的方法从路由获取

	propertyLock sync.RWMutex
	property     map[string]interface{} // 连接属性集合
}

func NewConnection(server msiface.IServer, conn *net.TCPConn, connID uint32, routers msiface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		MsgHandler: routers,
		property:   nil,
	}
	// 将conn加入connManager
	c.TcpServer.GetConnMgr().Add(c)

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

		// 根据head读data
		// TODO: msg.(*Message).Data是否改变接口变量?可能报错！
		// TODO: 判断dataLen为0的情况
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

		// 先判断是否开启工作池,用户忘记设置时,开goroutine处理业务
		if utils.GlobalConfig.WorkerPoolSize > 0 {
			c.MsgHandler.SendReqToTaskQueue(req)
		} else {
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
}

// StartWriter 开启写数据业务
func (c *Connection) StartWriter() {
	fmt.Println("writer goroutine is running...")
	defer fmt.Println(c.RemoteAddr().String(), "Writer exit!")

	// 循环等待Reader发的数据
	for {
		select {
		case data := <-c.msgChan:
			// channel中是[]byte类型，已pack
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("conn write data error:", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}
func (c *Connection) Start() {
	fmt.Println("Conn Start, ConnID is", c.ConnID)

	// 分别开启goroutine来read和write
	go c.StartReader()
	go c.StartWriter()

	// 钩子
	c.TcpServer.CallOnConnStart(c)
}

// Stop 关闭连接，同时通知Writer
func (c *Connection) Stop() {
	fmt.Println("Conn Stop, ConnID is", c.ConnID)

	// TODO：什么时候在调用stop前已经是true
	if c.isClosed {
		return
	}
	c.isClosed = true

	// 钩子
	c.TcpServer.CallOnConnStop(c)

	// 释放资源, close socket
	c.Conn.Close()
	c.TcpServer.GetConnMgr().Remove(c)
	close(c.ExitChan)
	close(c.msgChan)
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

// SendMsg 将[]byte的data=>封包=>TLV格式的[]byte类型的data
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
	c.msgChan <- msgByte

	return nil
}

// SetProperty 设置连接属性
func (c *Connection) SetProperty(key string, val interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if c.property == nil {
		c.property = make(map[string]interface{})
	}
	c.property[key] = val
}

// GetProperty 获取连接属性
func (c *Connection) GetProperty(key string) (val interface{}, err error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if v, ok := c.property[key]; ok {
		return v, nil
	}
	return nil, errors.New("property not found!")
}

// RemoveProperty 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
