package msnet

import (
	"fmt"
	"myserver/msiface"
	"myserver/utils"
	"net"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int

	MsgHandler msiface.IMsgHandler
	ConnMgr    msiface.IConnManager
	// 钩子函数
	OnConnStart func(conn msiface.IConnection) // 创建连接时
	OnConnStop  func(conn msiface.IConnection) // 销毁连接时
}

// NewServer 初始化Server模块
func NewServer(name string) msiface.IServer {
	s := &Server{
		Name:       utils.GlobalConfig.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalConfig.Host,
		Port:       utils.GlobalConfig.TcpPort,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server name is [%s], IP is %s, Port is %d, starting...\n",
		s.Name, s.IP, s.Port)

	// 开启goroutine去异步处理,在Serve()中阻塞
	go func() {
		// 开启工作池
		s.MsgHandler.StartWorkerPool()

		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp4 addr error: ", err)
			return
		}

		// 监听
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			panic(err)
		}
		fmt.Println("start server succeed, listening...")

		// 生成connID TODO: 自动生成ID的方法
		var cid uint32 = 0
		// 循环等待连接
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error:", err)
				continue
			}

			// 限制conn总数
			if s.ConnMgr.Count() >= utils.GlobalConfig.MaxConn {
				conn.Close()
				fmt.Println("error: too many conns!")
				continue
			}

			// TODO：生成connID
			// 将conn和处理业务的函数进行绑定=>得到新的Connection模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 开启连接
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// 关闭连接
	s.ConnMgr.CleanConn()
	fmt.Println("server stop(), clean conns completed")
}

func (s *Server) Serve() {
	s.Start()

	// TODO: 将启动服务独立出来，在这里可以扩展之后的其他业务

	// 阻塞,否则Start()的协程就结束了
	select {}
}

func (s *Server) AddRouter(msgID uint32, router msiface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("AddRouter succeed!")
}

func (s *Server) GetConnMgr() msiface.IConnManager {
	return s.ConnMgr
}

// SetOnConnStart 注册OnConnStart钩子方法
func (s *Server) SetOnConnStart(f func(msiface.IConnection)) {
	s.OnConnStart = f
}

// SetOnConnStop 注册OnConnStop钩子方法
func (s *Server) SetOnConnStop(f func(msiface.IConnection)) {
	s.OnConnStop = f
}

// CallOnConnStart 调用OnConnStart钩子方法
func (s *Server) CallOnConnStart(conn msiface.IConnection) {
	// 判断是否已注册钩子函数
	if s.OnConnStart != nil {
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 注册OnConnStop钩子方法
func (s *Server) CallOnConnStop(conn msiface.IConnection) {
	// 判断是否已注册钩子函数
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	}
}
