package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CorsFilter() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           86400,
	})
}
