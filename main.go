package main

import (
	"commette-chat/config"
	"commette-chat/routes"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))
	config.ConnectDB()
	defer config.DisconnectDB()
	app.Static("/", "./public")
	routes.SetupRoutes(app)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Cerrando la aplicación...")
		if err := app.Shutdown(); err != nil {
			log.Fatalf("Error cerrando la aplicación: %v", err)
		}
	}()

	if err := app.Listen(":4000"); err != nil {
		log.Fatalf("Error iniciando el servidor: %v", err)
	}
}
