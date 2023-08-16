package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/Eberewill/emotechat/initializers"
	"github.com/Eberewill/emotechat/routes"
	"github.com/Eberewill/emotechat/services"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

//utils

type MessageObj struct {
	Message  string `json:"message"`
	Sender   uint   `json:"sender"`
	Reciever uint   `json:"receiver"`
	Time     string `json:"time"`
	Id       int    `json:"id"`
}

func ConvertJSONToObject(jsonStr []byte) (MessageObj, error) {
	var obj MessageObj
	err := json.Unmarshal([]byte(jsonStr), &obj)
	if err != nil {
		return MessageObj{}, err
	}
	return obj, nil
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

			obj, err := ConvertJSONToObject(msg)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			fmt.Println("obj.Sender", obj.Sender)
			fmt.Println("obj.Reciever", obj.Reciever)
			// Save the message to the database
			err = services.SaveMessage(obj.Sender, obj.Reciever, string(obj.Message))
			if err != nil {
				log.Println("message service error:", err)
			}

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
