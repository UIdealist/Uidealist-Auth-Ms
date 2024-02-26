package routes

import (
	"idealist/app/controllers"

	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for POST method:
	route.Post("/user/sign/in", controllers.UserSignIn) // auth, return Access & Refresh tokens
	route.Post("/user/sign/up", controllers.UserSignUp) // create a new user
}
