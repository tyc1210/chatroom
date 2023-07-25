package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":20001")
	if err != nil {
		panic(err)
	}

	done := make(chan struct{})
	go func() {
		// 将从 conn（一个实现了 net.Conn 接口的对象）读取的数据复制到标准输出
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()
	// 从标准输入（os.Stdin）读取的数据复制到连接（conn）中
	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
