package main

import (
	"os"

	"idealist/pkg/configs"
	"idealist/pkg/middleware"
	"idealist/pkg/routes"
	"idealist/pkg/utils"
	"idealist/platform/database"

	"github.com/gofiber/fiber/v2"

	_ "idealist/docs" // load API Docs files (Swagger)

	_ "github.com/joho/godotenv/autoload" // load .env file automatically
)

// @Idealist API
// @version 1.0
// @description Idealist project API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email edgardanielgd123@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Define Fiber config.
	config := configs.FiberConfig()

	// Create database connection.
	database.OpenDBConnection()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.

	// Routes.
	routes.SwaggerRoute(app)  // Register a route for API Docs (Swagger).
	routes.PublicRoutes(app)  // Register a public routes for app.
	routes.PrivateRoutes(app) // Register a private routes for app.
	routes.NotFoundRoute(app) // Register route for 404 Error.

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
