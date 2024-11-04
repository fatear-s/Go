package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	b, err := net.Listen("tcp", "127.0.0.1:1080")
	if err != nil {
		log.Printf("监听失败:%w", err)
	}
	conn, err := b.Accept()
	if err != nil {
		log.Printf("连接失败:%w", err)
	}
	process(conn)
}
func process(conn net.Conn) {
	defer conn.Close()
	render := bufio.NewReader(conn)
	printMsg := []byte(fmt.Sprintf("%s connection\n>", conn.RemoteAddr()))
	sendMsg := []byte(fmt.Sprintf("%s connection sucessful\n", conn.LocalAddr()))
	_, err := conn.Write(sendMsg)
	if err != nil {
		log.Printf("send error:%w", err)
	}
	log.Println(string(printMsg))
	for {
		buf, err := render.ReadByte()
		if err != nil {
			break
		}
		_, err = conn.Write([]byte{buf})
		if err != nil {
			log.Printf("send error:%w", err)
		}
		log.Printf(string(buf))
	}
}
