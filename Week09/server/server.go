package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Conn struct {
	ip       string
	Conn     net.Conn
	MsgChan  chan []byte
	ExitChan chan bool
	Closed   bool
}

func main() {

	c := &Conn{
		ip:       "127.0.0.1:8000",
		MsgChan:  make(chan []byte),
		ExitChan: make(chan bool),
	}
	listen, err := net.Listen("tcp", c.ip)
	if err != nil {
		log.Fatal("listen error:%v\n", err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal("accept error:%v\n", err)
			continue
		}
		c.Conn = conn
		go c.readConn()
		go c.writeConn()
	}
}
func (c *Conn) readConn() {
	defer c.Stop()
	rd := bufio.NewReader(c.Conn)

	for {
		line, _, err := rd.ReadLine()
		fmt.Println(string(line))
		if err != nil {
			log.Printf("read error:%v\n", err)
			return
		}
		c.MsgChan <- []byte("hello " + string(line))
	}
}
func (c *Conn) writeConn() {
	defer c.Stop()
	wt := bufio.NewWriter(c.Conn)
	for {
		select {
		case data := <-c.MsgChan:
			wt.Write(data)
			wt.Flush()
		case <-c.ExitChan:
			return
		}
	}
}
func (c *Conn) Stop() {
	if c.Closed {
		return
	}
	c.Closed = true
	c.ExitChan <- true
	c.Conn.Close()
	close(c.ExitChan)
	close(c.MsgChan)
}
