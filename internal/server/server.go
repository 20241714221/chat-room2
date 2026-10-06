package server

import (
	"fmt"
	"net"
	"chat-room/internal/model"

)


type Server struct {
	Addr string
	Listener net.Listener
	Room *Room
}

func NewServer(addr string) *Server {
	return &Server{
		Addr: addr,
		Room:NewRoom(),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen(
		"tcp",
		s.Addr,
	)
	if err != nil {
		return err
	}
	s.Listener = listener
	go s.Room.Run()
	fmt.Println("服务器启动，监听:", s.Addr)
	return nil
}

func (s *Server) Accept() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println("连接失败:",err) 
			continue
		}
		client := &Client{
			ID:GenerateID(),
			Conn:conn,
			Room:s.Room,
			Send:make(chan model.Message,100),
		}
		go client.Handle()
	}
}