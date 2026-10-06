package main

import (
	"fmt"

	"chat-room/internal/server"
)

func main(){
	s := server.NewServer(":8080")
	err := s.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	s.Accept()
}