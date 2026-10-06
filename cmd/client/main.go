package main

import (
    "fmt"
    "os"

    "chat-room/internal/client"
)

func main() {
    addr := os.Getenv("SERVER_ADDR")
    if addr == "" {
        addr = "localhost:8080"
    }

    c, err := client.NewClient(addr)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer c.Conn.Close()

    fmt.Print("请输入昵称:")
    name, _ := c.Input.ReadString('\n')

    c.Conn.Write([]byte(name))

    fmt.Println("进入聊天室")

    go c.Receive()
    c.Send()
}