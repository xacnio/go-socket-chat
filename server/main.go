package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Client struct {
	Name string
	Conn net.Conn
}

type Server struct {
	Clients []*Client
	Mu      sync.Mutex
}

type Message struct {
	Client *Client
	Message string
}

var (
	Sw Server
	MsgChannel chan Message
)

func ListenChannel() {
	for {
		msg := <-MsgChannel
		for i := 0; i < len(Sw.Clients); i++ {
			a := Sw.Clients[i]
			if a.Conn == msg.Client.Conn {
				continue
			}
			if _, err := a.Conn.Write([]byte(fmt.Sprintf("[%s]: << [%s]", msg.Client.Name, msg.Message))); err == io.EOF {
				Sw.Mu.Lock()
				Sw.Clients = append(Sw.Clients[:i], Sw.Clients[i+1:]...)
				i--
				a.Conn.Close()
				Sw.Mu.Unlock()
			}
		}
	}
}

func main() {
	MsgChannel = make(chan Message)
	Sw.Clients = make([]*Client, 0)
	ln, err := net.Listen("tcp", ":3500")
	if err != nil {
		panic(err)
	}
	go ListenChannel()
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			read := make([]byte, 64)
			n, err := conn.Read(read)
			if err != nil {
				conn.Close()
			} else {
				Sw.Mu.Lock()
				client := &Client{Name: string(read)[0:n], Conn: conn}
				Sw.Clients = append(Sw.Clients, client)
				fmt.Printf("%s (%s) connected!\n", client.Name, client.Conn.RemoteAddr())
				Sw.Mu.Unlock()
				SetupClient(client)
			}
		}()
	}
}

func SetupClient(c *Client) {
	for {
		read := make([]byte, 512)
		n, err := c.Conn.Read(read)
		if err != nil {
			fmt.Printf("%s (%s) disconnected!\n", c.Name, c.Conn.RemoteAddr())
			DeleteClient(c)
			return
		}
		msg := string(read)[0:n]
		fmt.Printf("[%s]: %s\n", c.Name, msg)
		MsgChannel<-Message{Message: msg, Client: c}
	}
}

func DeleteClient(c *Client) {
	Sw.Mu.Lock()
	defer Sw.Mu.Unlock()
	for i := 0; i < len(Sw.Clients); i++ {
		if Sw.Clients[i] == c {
			Sw.Clients = append(Sw.Clients[:i], Sw.Clients[i+1:]...)
			i--
			c.Conn.Close()
		}
	}
}
