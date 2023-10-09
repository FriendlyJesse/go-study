package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// 添加用户的函数
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	// 启动监听当前 user channel 消息的 goroutine
	go user.ListenMessage()

	return user
}

// 监听 User channel 的方法，有消息就发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		// 发送消息，它需要二进制数组，所以做一个转化
		this.conn.Write([]byte(msg + "\n"))
	}
}
