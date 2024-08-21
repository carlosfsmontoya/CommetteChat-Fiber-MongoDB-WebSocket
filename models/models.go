package models

import "time"

// User representa un usuario en el chat
type User struct {
	IDUser int `json:"id_user"`
}

// Conversation representa una conversación en el chat
type Conversation struct {
	IDParticipants []string `json:"id_participants"`
}

// Message representa un mensaje en el chat
type Message struct {
	ConversationID string    `json:"conversation_id" bson:"conversation_id"`
	SenderID       string    `json:"sender_id" bson:"sender_id"`
	Content        string    `json:"content" bson:"content"`
	Timestamp      time.Time `json:"timestamp" bson:"timestamp"`
}
