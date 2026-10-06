package server

import (
	"chat-room/internal/model"
	"fmt"
	"time"
)

type Room struct {
	// 在线用户
	clients map[int64]*Client
	// 用户加入
	join chan *Client
	// 用户退出
	leave chan *Client
	// 广播消息
	broadcast chan model.Message
	// 查询用户列表
	listUsers chan *Client
}

func NewRoom() *Room {
	return &Room{
		clients:   make(map[int64]*Client),
		join:      make(chan *Client),
		leave:     make(chan *Client),
		broadcast: make(chan model.Message, 100),
		listUsers: make(chan *Client),
	}
}

// CheckName 检查用户名是否合法
// 返回错误信息，如果合法返回空字符串
func (r *Room) CheckName(name string) string {
	// 检查是否为空
	if name == "" {
		return "昵称不能为空"
	}
	// 检查长度
	if len(name) > 10 {
		return "昵称长度不能超过10个字符"
	}
	// 检查重名
	for _, c := range r.clients {
		if c.Name == name {
			return "昵称已经存在"
		}
	}
	return ""
}

func (r *Room) Run() {
	for {
		select {
		// 查询用户列表
		case client := <-r.listUsers:
			var userList string
			count := 1
			for _, c := range r.clients {
				userList += fmt.Sprintf("%d. %s\n", count, c.Name)
				count++
			}
			if userList == "" {
				userList = "当前在线用户: 无\n"
			} else {
				userList = "当前在线用户:\n" + userList
			}
			msg := model.Message{
				Type:   "list",
				Sender: "系统",
				Content: userList,
				Time:   time.Now().Unix(),
			}
			client.Send <- msg
		// 用户加入
		case client := <-r.join:
			// 保存用户
			r.clients[client.ID] = client
			// 当前在线人数
			count := len(r.clients)
			// 创建系统消息
			msg := model.Message{
				Type:   "join",
				Sender: "系统",
				Content: fmt.Sprintf(
					"%s 加入聊天室 当前在线人数:%d",
					client.Name,
					count,
				),
				Time: time.Now().Unix(),
			}
			// 放入广播队列
			r.broadcast <- msg
		// 用户退出
		case client := <-r.leave:
			// 删除用户
			delete(
				r.clients,
				client.ID,
			)
			// 关闭发送channel
			close(client.Send)
			count := len(r.clients)
			msg := model.Message{
				Type:   "leave",
				Sender: "系统",
				Content: fmt.Sprintf(
					"%s 离开聊天室 当前在线人数:%d",
					client.Name,
					count,
				),
				Time: time.Now().Unix(),
			}
			r.broadcast <- msg
		// 广播消息
		case msg := <-r.broadcast:
			for _, client := range r.clients {
				select {
				case client.Send <- msg:
				default:
					fmt.Println(
						client.Name,
						"消息队列已满",
					)
				}
			}
		}
	}
}
