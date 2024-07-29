package controllers

import (
	"commette-chat/models"
	"log"
	"time"

	"github.com/gofiber/websocket/v2"
)

func HandleWebSocket(c *websocket.Conn) {
	defer func() {
		mutex.Lock()
		delete(clients, c)
		mutex.Unlock()
		c.Close()
	}()

	mutex.Lock()
	clients[c] = true
	mutex.Unlock()

	for {
		var msg models.Message
		if err := c.ReadJSON(&msg); err != nil {
			log.Printf("error: %v", err)
			break
		}
		msg.Timestamp = time.Now().Format("2006-01-02T15:04:05")
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client := range clients {
			if err := client.WriteJSON(msg); err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

func init() {
	go handleMessages()
}
