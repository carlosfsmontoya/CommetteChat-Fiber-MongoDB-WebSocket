package controllers

import (
	"commette-chat/config"
	"commette-chat/models"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	userCollection := config.GetCollection(config.DB, "users")

	// Verificar si el id_user ya existe
	var existingUser models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"id_user": user.IDUser}).Decode(&existingUser)
	if err == nil {
		// Si no hay error, significa que el usuario ya existe
		return c.Status(fiber.StatusConflict).SendString("User already exists")
	}

	// Insertar el nuevo usuario si no existe
	newUser := bson.M{
		"id_user": user.IDUser,
	}
	result, err := userCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}
