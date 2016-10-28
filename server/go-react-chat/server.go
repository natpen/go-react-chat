package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/twinj/uuid"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

// GLOBALS

const dbUrlEnvKey string = "DATABASE_URL"
const messagesPerPage int = 8

var db *sql.DB
var chat ChatRoom
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // not checking origin
}

func staticFiles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./client/"+r.URL.Path)
}

// this is also the handler for joining to the chat
func wsHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Error upgrading to websocket:", err)
		return
	}
	go func() {
		// first message has to be the name
		_, message, err := conn.ReadMessage()
		client := chat.Join(string(message), conn)
		if client == nil || err != nil {
			conn.Close() // closing connection to indicate failed Join
			return
		}

		// then watch for incoming messages
		for {
			_, message, err := conn.ReadMessage()
			if err != nil { // if error then assuming that the connection is closed
				client.Exit()
				return
			}

			var messageUnmarshaled map[string]interface{}
			json.Unmarshal(message, &messageUnmarshaled)

			unixTimestampSec := messageUnmarshaled["timestamp"].(float64) / 1000.0
			unixTimestampNSec := (unixTimestampSec - math.Floor(unixTimestampSec)) * 1000000000.0
			timestamp := time.Unix(int64(unixTimestampSec), int64(unixTimestampNSec))

			mId, err := uuid.Parse(fmt.Sprintf("%s", messageUnmarshaled["id"]))
			if err != nil {
				fmt.Println(err)
				client.Send([]Message{Message{uuid.NewV4(), "system-message", "", time.Now(), "Error sending message"}})
				return
			}

			newMessage := Message{
				Id:        mId,
				Type:      messageUnmarshaled["type"].(string),
				Name:      client.name,
				Timestamp: timestamp,
				Text:      messageUnmarshaled["text"].(string),
			}

			client.NewMessage(newMessage)
		}

	}()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err = sql.Open("postgres", os.Getenv(dbUrlEnvKey))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	uuid.Init()

	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", staticFiles)
	chat.Init()
	fmt.Println("\nSuccess! Please navigate your browser to http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}
