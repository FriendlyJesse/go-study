package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	// 在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	// 消息广播的 channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// 监听消息，有消息就将消息发送给在线用户
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message
		// fmt.Println(msg)
		// 将消息发送给所有的在线用户
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播消息的方法
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ": " + msg
	// 发送消息给用户
	this.Message <- sendMsg
}

// 处理业务
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("用户建立链接")

	user := NewUser(conn)

	// 将用户加入 OnlineMap 中
	userAddr := conn.RemoteAddr().String()
	this.mapLock.Lock()
	this.OnlineMap[userAddr] = user
	this.mapLock.Unlock()

	// 广播用户上线消息
	this.BroadCast(user, "已上线！")
}

// 启动服务
func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	defer listener.Close()
	fmt.Printf("服务已启动：%s:%d\n", this.Ip, this.Port)

	// 启动监听用户上线 message 的 goroutine
	go this.ListenMessager()

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
