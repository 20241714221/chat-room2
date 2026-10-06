package model


type Message struct {
	Type string `json:"type"`
	Sender string `json:"sender"`
	Content string `json:"content"`
	Time int64 `json:"time"`
}