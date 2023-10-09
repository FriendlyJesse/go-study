package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}

	return server
}

// 处理业务
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("用户建立链接")
}

// 启动服务
func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	defer listener.Close()
	fmt.Printf("服务已启动：%s:%d", this.Ip, this.Port)

	for {
		// 等待并返回一个已连接的侦听器
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listener accept err:", err)
			continue
		}

		go this.Handler(conn)
	}
}
