package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CorsConfig() func(*fiber.Ctx) error {
	// frontendURL := os.Getenv("FRONTEND_URL_DEV")
	frontendURL := os.Getenv("FRONTEND_URL")
	return cors.New(cors.Config{
		// AllowOrigins: frontendURL,
		AllowOrigins:     "http://localhost:3000," + frontendURL,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		MaxAge:           600,
	})

}
