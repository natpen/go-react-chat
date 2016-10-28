package main

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	name     string
	conn     *websocket.Conn
	chatRoom *ChatRoom
}

// Client has a new message to broadcast
func (client *Client) NewMessage(message Message) {
	client.chatRoom.AddMessage(message)
	AddUserOrUpdateLastActive(client.name)
}

// Exiting out
func (client *Client) Exit() {
	client.chatRoom.Leave(client.name)
}

// Sending message block to the client
func (client *Client) Send(messages []Message) {
	client.conn.WriteJSON(messages)
}
