package routes

import (
	"github.com/Eberewill/emotechat/middleware"
	"github.com/Eberewill/emotechat/services"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	app.Post("/api/auth/register", services.Register)
	app.Post("/api/auth/login", services.Login)
	app.Post("/api/auth/logout", middleware.RequireAuth, services.Logout)
	app.Get("/api/auth/validate", middleware.RequireAuth, services.Validate)

}
