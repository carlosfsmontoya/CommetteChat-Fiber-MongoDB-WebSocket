package routes

import (
	"commette-chat/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api")

	api.Post("/users", controllers.InsertUser)
	api.Post("/conversations", controllers.StartConversation)
	api.Post("/messages", controllers.InsertMessage)

	app.Get("/ws/:conversationID", websocket.New(controllers.HandleWebSocket))

}
