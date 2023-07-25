package main

import (
	"bufio"
	"chatroom/model"
	"chatroom/util"
	"fmt"
	"log"
	"net"
)

var (
	// 存放进入游戏消息的channel
	enterChannel = make(chan *model.User)
	// 存放退出游戏消息的channel
	leaveChannel = make(chan *model.User)
	// 存放广播消息的channel
	msgChannel = make(chan model.Message, 8)

	address = ":20001"
)

func main() {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("start server address:%s", address)

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		go handleConn(conn)
	}

}

func handleConn(conn net.Conn) {
	defer conn.Close()
	// 构建用户 为用户分配 channel
	user := &model.User{
		Id:             util.GetId(),
		Addr:           conn.RemoteAddr().String(),
		MessageChannel: make(chan model.Message, 10),
	}
	go user.StartListenAndSend(conn)
	// 向用户发送欢迎消息
	user.MessageChannel <- model.Message{Uid: "群消息", Msg: "欢迎加入三年二班群"}
	// 发送群消息
	msgChannel <- model.Message{Uid: "群消息", Msg: fmt.Sprintf("用户%s加入了群聊", user.Id)}
	// 将用户加入管道发送消息
	enterChannel <- user
	// 循环读取用户输入
	input := bufio.NewScanner(conn)
	for input.Scan() {
		msg := model.Message{Uid: user.Id, Msg: input.Text()}
		msgChannel <- msg
	}
	if err := input.Err(); err != nil {
		log.Println("读取错误：", err)
	}
	// 用户离开
	leaveChannel <- user
	msgChannel <- model.Message{Uid: user.Id, Msg: fmt.Sprintf("用户%s离开了群聊", user.Id)}
}

func broadcaster() {
	// map 保存所有群聊用户信息
	users := make(map[string]*model.User)
	for true {
		select {
		case user := <-enterChannel:
			users[user.Id] = user
		case user := <-leaveChannel:
			delete(users, user.Id)
			close(user.MessageChannel)
		case msg := <-msgChannel:
			for _, user := range users {
				if user.Id != msg.Uid {
					user.MessageChannel <- msg
				}
			}
		}
	}
}
