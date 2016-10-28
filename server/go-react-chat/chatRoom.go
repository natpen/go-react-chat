package main

import (
	"github.com/gorilla/websocket"
	"github.com/twinj/uuid"
	"sync"
	"time"
)

type ChatRoom struct {
	clients    map[string]Client
	clientsMtx sync.Mutex
	queue      chan Message
}

// Initializing the chatroom
func (chatRoom *ChatRoom) Init() {
	chatRoom.queue = make(chan Message, 5)
	chatRoom.clients = make(map[string]Client)

	go func() {
		for {
			chatRoom.BroadCast()
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

// Registering a new client
// returns pointer to a Client, or Nil, if the name is already taken
func (chatRoom *ChatRoom) Join(name string, conn *websocket.Conn) *Client {
	defer chatRoom.clientsMtx.Unlock()

	chatRoom.clientsMtx.Lock()
	if _, exists := chatRoom.clients[name]; exists && len(name) >= 3 {
		return nil
	}
	client := Client{
		name:     name,
		conn:     conn,
		chatRoom: chatRoom,
	}
	chatRoom.clients[name] = client

	lastActive, _ := GetUserLastActive(name)

	client.Send([]Message{Message{uuid.NewV4(), "client-handshake", "", lastActive, ""}})

	messages := GetMessagesForUser(client.name, time.Now())

	AddUserOrUpdateLastActive(name)

	client.Send(messages)

	chatRoom.AddMessage(Message{uuid.NewV4(), "system-message", "", time.Now(), name + " has joined the chat."})
	return &client
}

// Leaving the chatroom
func (chatRoom *ChatRoom) Leave(name string) {
	chatRoom.clientsMtx.Lock()
	delete(chatRoom.clients, name)
	AddUserOrUpdateLastActive(name)
	chatRoom.clientsMtx.Unlock()
	chatRoom.AddMessage(Message{uuid.NewV4(), "system-message", "", time.Now(), name + " has left the chat."})
}

// Adding message to queue
func (chatRoom *ChatRoom) AddMessage(message Message) {
	chatRoom.queue <- message
	if message.Name != "" {
		StoreMessage(message)
	}
}

// Broadcasting all the messages in the queue in one block
func (chatRoom *ChatRoom) BroadCast() {

	messages := make([]Message, 0)

infLoop:
	for {
		select {
		case message := <-chatRoom.queue:
			messages = append(messages, message)
		default:
			break infLoop
		}
	}
	if len(messages) > 0 {
		for _, client := range chatRoom.clients {
			client.Send(messages)
		}
	}
}
