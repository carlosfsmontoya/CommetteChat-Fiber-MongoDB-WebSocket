package models

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
	ConversationID string `json:"conversation_id"`
	SenderID       string `json:"sender_id"`
	Content        string `json:"content"`
	Timestamp      string `json:"timestamp"`
}
