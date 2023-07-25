package model

import (
	"fmt"
	"net"
)

type User struct {
	Id             string
	Addr           string
	MessageChannel chan Message
}

// StartListenAndSend 从管道中读取数据发送到客户端
func (u User) StartListenAndSend(conn net.Conn) {
	for msg := range u.MessageChannel {
		fmt.Fprintln(conn, msg.Uid+":"+msg.Msg)
	}
}
