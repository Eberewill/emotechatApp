package routes

import (
	"github.com/Eberewill/emotechat/middleware"
	"github.com/Eberewill/emotechat/services"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", services.RegisterUser)
	auth.Post("/login", services.LoginUser)
	auth.Post("/logout", middleware.RequireAuth, services.LogoutUser)
	auth.Get("/validate", middleware.RequireAuth, services.ValidateUser)

	profile := api.Group("/profile", middleware.RequireAuth)
	profile.Get("/users", services.GetAllUsers)

	conversations := api.Group("/conversations", middleware.RequireAuth)
	conversations.Get("/", services.FetchConversations)
	// conversations.Get("/test/message", services.TestSaveMessage)
	// conversations.Get("/test/message/all", services.TestFetchConversations)

	password := api.Group("/password")
	password.Post("/reset-request", services.RequestPasswordReset)
	password.Post("/reset", services.ResetPassword)
}
