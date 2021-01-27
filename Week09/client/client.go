package main

import (
	"fmt"
	"log"
	"net"
)

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatal("connect failed:%v", err)
	}
	defer conn.Close()
	msgClient := "request connection!\n"
	sendN, err := conn.Write([]byte(msgClient))
	fmt.Println("send", sendN)
	if err != nil {
		log.Fatal("send msg failed:%v\n", err)
	}

	buf := make([]byte, 1024)
	respN, err := conn.Read(buf)
	fmt.Println("resp", respN)

	if err != nil {
		log.Fatal("recevid msg failed:%v\n", err)
	}
	fmt.Println(string(buf))
}
