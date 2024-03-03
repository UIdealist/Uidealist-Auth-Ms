package routes

import (
	"idealist/app/controllers"
	"idealist/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for POST method:
	route.Post(
		"/user/sign/out", middleware.JWTProtected(),
		middleware.JWTExpirationChecker(), controllers.UserSignOut) // de-authorization user
	route.Post("/token/renew", middleware.JWTProtected(), controllers.RenewTokens) // renew Access & Refresh tokens
}
