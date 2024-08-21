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

	// Imprime el cuerpo de la solicitud
	body := c.Body()
	log.Println("Received body:", string(body))

	// Analiza el cuerpo de la solicitud en el modelo Message
	if err := c.BodyParser(&message); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Imprime el mensaje analizado
	log.Printf("Parsed message: %+v\n", message)

	// Analiza el campo timestamp
	parsedTime, err := time.Parse(time.RFC3339, message.Timestamp.Format(time.RFC3339))
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Invalid timestamp format")
	}

	// Obtén la colección de mensajes
	msgCollection := config.GetCollection(config.DB, "messages")

	// Crea un nuevo documento de mensaje
	newMessage := bson.M{
		"conversation_id": message.ConversationID,
		"sender_id":       message.SenderID,
		"content":         message.Content,
		"timestamp":       parsedTime,
	}

	// Inserta el nuevo mensaje en la base de datos
	result, err := msgCollection.InsertOne(context.TODO(), newMessage)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to insert message")
	}

	// Responde con el ID del nuevo mensaje insertado
	return c.Status(fiber.StatusOK).JSON(result)
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
