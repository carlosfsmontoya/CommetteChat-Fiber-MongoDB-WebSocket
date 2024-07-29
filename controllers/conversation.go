package controllers

import (
	"commette-chat/config"
	"commette-chat/models"
	"context"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"go.mongodb.org/mongo-driver/bson"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan models.Message)
var mutex = &sync.Mutex{}

func InsertMessage(c *fiber.Ctx) error {
	var message models.Message

	if err := c.BodyParser(&message); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	log.Print(message)

	parsedTime, err := time.Parse("2006-01-02T15:04", message.Timestamp)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid timestamp format")
	}
	msgCollection := config.GetCollection(config.DB, "messages")
	newMessage := bson.M{
		"conversation_id": message.ConversationID,
		"sender_id":       message.SenderID,
		"content":         message.Content,
		"timestamp":       parsedTime,
	}

	result, err := msgCollection.InsertOne(context.TODO(), newMessage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

func StartConversation(c *fiber.Ctx) error {
	var conversation models.Conversation
	if err := c.BodyParser(&conversation); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	participantIDs := conversation.IDParticipants

	convCollection := config.GetCollection(config.DB, "conversations")
	newConversation := bson.M{
		"id_participants": participantIDs,
		"id_last_message": nil,
	}
	result, err := convCollection.InsertOne(context.TODO(), newConversation)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}
