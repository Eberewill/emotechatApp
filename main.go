package main

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/Eberewill/emotechat/initializers"
	"github.com/Eberewill/emotechat/routes"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	// Define a list of allowed origins
	allowedOrigins := []string{
		"http://127.0.0.1:5173",
		"*",
		// Add more allowed origins as needed
	}

	// Join the allowed origins into a comma-separated string
	allowOriginsString := strings.Join(allowedOrigins, ",")

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     allowOriginsString,
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	// WebSocket upgrade middleware
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	var (
		clients     = make(map[*websocket.Conn]bool)
		clientsLock sync.Mutex
	)

	// WebSocket endpoint
	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		clientsLock.Lock()
		clients[c] = true
		clientsLock.Unlock()

		defer func() {
			clientsLock.Lock()
			delete(clients, c)
			clientsLock.Unlock()
			c.Close()
		}()

		log.Println("Client connected")

		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			log.Printf("recv: %s", msg)

			// Save the message to the database
			//err = message.SaveMessage(senderID, receiverID, string(msg))
			//	if err != nil {
			//		log.Println("message service error:", err)
			//	}

			// Broadcast the message to other clients
			clientsLock.Lock()
			for client := range clients {
				if err := client.WriteMessage(mt, msg); err != nil {
					log.Println("write:", err)
				}
			}
			clientsLock.Unlock()
		}
	}))

	// Start the Fiber app

	routes.Setup(app)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	err := app.Listen(":" + port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
