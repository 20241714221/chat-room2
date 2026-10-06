package client

import (
	"bufio"
	"fmt"
	"net"
	"encoding/json"
	"chat-room/internal/model"
	"os"
	"strings"
	"time"
)
type Client struct {
	Conn net.Conn
	Input *bufio.Reader
	lastTime int64 // 上次消息时间
}
func NewClient(addr string) (*Client,error){
	conn,err:=net.Dial(
		"tcp",
		addr,
	)
	if err!=nil{
		return nil,err
	}
	return &Client{
		Conn:conn,
		Input:bufio.NewReader(os.Stdin),
		lastTime:0,
	},nil
}
func (c *Client) Send(){
	reader:=bufio.NewReader(
		os.Stdin,
	)
	for{
		msg,err:=reader.ReadString('\n')
		if err!=nil{
			return
		}
		// 去除换行符
		msg = strings.TrimRight(msg, "\n")
		// 空消息不发送
		if msg == "" {
			continue
		}
		// 检查是否是quit命令
		if msg == "/quit" {
			c.Conn.Write([]byte(msg + "\n"))
			return
		}
		_,err=c.Conn.Write(
			[]byte(msg + "\n"),
		)
		if err!=nil{
			return
		}
	}
}
func(c *Client) Receive(){
	reader:=bufio.NewReader(c.Conn)
	for{
		data,err:=reader.ReadBytes('\n')
		if err!=nil{
			return
		}
		var msg model.Message
		err=json.Unmarshal(
			data,
			&msg,
		)//解析之后放在msg
		if err!=nil{
			continue
		}
		switch msg.Type{
		case "error":
			fmt.Println(
				"[错误]",
				msg.Content,
			)
		case "chat":
			// 判断是否需要显示时间（超过60秒）
			showTime := c.lastTime == 0 || (msg.Time - c.lastTime) > 10
			if showTime {
				fmt.Printf("[%s] [%s]: %s\n",
					time.Unix(msg.Time, 0).Format("15:04:05"),
					msg.Sender,
					msg.Content,
				)
			} else {
				fmt.Printf("[%s]: %s\n",
					msg.Sender,
					msg.Content,
				)
			}
			c.lastTime = msg.Time
		case "join":
			fmt.Println(
				"[系统]",
				msg.Content,
			)
		case "leave":
			fmt.Println(
				"[系统]",
				msg.Content,
			)
		case "list":
			fmt.Print(
				msg.Content,
			)
		case "quit":
			fmt.Println(
				"[系统]",
				msg.Content,
			)
			// 收到quit消息，关闭连接让Send也退出
			c.Conn.Close()
			return
		}
	}
}