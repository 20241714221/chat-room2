package server

import (
	"bufio"
	"chat-room/internal/model"
	"encoding/json"
	"net"
	"strings"
	"time"
)

type Client struct {
	ID   int64
	Name string
	Conn net.Conn
	Room *Room
	Send chan model.Message
}

func (c *Client) Write() {
	for msg := range c.Send { //一直等待：
		data, err := json.Marshal(msg)//Go结构体转json字符串
		if err != nil {
			continue
		}
		_, err = c.Conn.Write(
			append(data, '\n'),
		)
		if err != nil {
			return
		}
	}
}

func (c *Client) Handle() {
	defer c.Conn.Close()
	defer func() {
		c.Room.leave <- c
	}()
	reader := bufio.NewReader(c.Conn)

	// 用户名校验循环
	for {
		name, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		name = strings.TrimSpace(name)

		// 校验用户名
		errMsg := c.Room.CheckName(name)
		if errMsg != "" {
			// 发送错误消息给客户端
			errorMsg := model.Message{
				Type:    "error",
				Content: errMsg,
			}
			data, err := json.Marshal(errorMsg)//转换成字节数组
			if err == nil {
				c.Conn.Write(append(data, '\n'))
			}
			continue
		}

		c.Name = name
		break
	}

	go c.Write()
	c.Room.join <- c
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		msg = strings.TrimSpace(msg)

		// 处理命令
		switch msg {
		case "/list":
			c.Room.listUsers <- c
		case "/quit":
			// 发送退出系统消息
			quitMsg := model.Message{
				Type:    "quit",
				Sender:  c.Name,
				Content: c.Name + "退出聊天室",
				Time:    time.Now().Unix(),
			}
			c.Room.broadcast <- quitMsg
			return
		default:
			if msg!= ""{
			message := model.Message{
				Type:    "chat",
				Sender:  c.Name,
				Content: msg,
				Time:    time.Now().Unix(),
			}
			c.Room.broadcast <- message
		}
		}
	}
}
