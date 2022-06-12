package utils

import "myserver/msiface"

var GlobalConfig *GlobalCfg

type GlobalCfg struct {
	/***	Server config	***/
	TcpServer msiface.IServer // 全局Server对象
	Host string
	Name string
	TcpPort int // 监听的端口号

	/***	App Config	***/
	Version string
	MaxConn int // 最大连接数
	MaxPacketSize uint32 // 数据包的最大值
}

// 自动初始化默认配置
func init() {
	// 默认配置
	GlobalConfig = &GlobalCfg{
		Name: "myserver app",
		Version: "v0.1",
		TcpPort: 7777,
		Host: "0.0.0.0", //本地全部IP
		MaxConn: 100,
		MaxPacketSize: 4096,
	}

	// 加载用户的配置
	GlobalConfig.Reload()
}

func (g GlobalCfg) Reload() {
	// TODO: 使用viper加载配置
}