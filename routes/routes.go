package routes

import (
	"github.com/Eberewill/emotechat/middleware"
	"github.com/Eberewill/emotechat/services"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	app.Post("/api/auth/register", services.RegisterUser)
	app.Post("/api/auth/login", services.LoginUser)
	app.Post("/api/auth/logout", middleware.RequireAuth, services.LogoutUser)
	app.Get("/api/auth/validate", middleware.RequireAuth, services.ValidateUser)
	app.Get("/api/test/message", services.TestSaveMessage)
	app.Get("/api/test/message/all", services.TestFetchConversations)

}
