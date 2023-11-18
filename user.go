package main

import "net"

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// 创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	// 启动监听当前user channel消息的goroutine
	go user.ListenMessage()

	return user
}

// Online 用户上线
func (this *User) Online() {
	// 用户上线
	// 1. 将用户加入onlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	// 2. 广播当前用户上线消息
	this.server.Broadcast(this, "已经上线")
}

// Offline 用户下线
func (this *User) Offline() {
	// 用户下线
	// 1. 将用户从onlineMap中删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	// 2. 广播当前用户上线消息
	this.server.Broadcast(this, "下线")
}

// DoMessage 用户处理消息
func (this *User) DoMessage(msg string) {
	this.server.Broadcast(this, msg)
}

// 监听当前User channel的方法，一旦有消息，就直接发给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
