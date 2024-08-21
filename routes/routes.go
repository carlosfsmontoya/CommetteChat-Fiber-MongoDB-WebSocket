package routes

import (
	"commette-chat/controllers"
	"commette-chat/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api")

	api.Post("/users", middleware.SecretKeyRequired(), controllers.InsertUser)

	api.Post("/conversations", controllers.StartConversation)
	api.Post("/messages", controllers.InsertMessage)

	api.Get("/hello", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/ws/:conversationID", websocket.New(controllers.HandleWebSocket))

}
