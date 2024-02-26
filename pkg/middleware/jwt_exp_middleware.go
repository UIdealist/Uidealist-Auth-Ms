package middleware

import (
	"idealist/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func JWTExpirationChecker() func(*fiber.Ctx) error {

	return func(c *fiber.Ctx) error {
		now := time.Now().Unix()

		claims, err := utils.ExtractTokenMetadata(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		expires := claims.Expires

		if now > expires {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   "unauthorized, check expiration time of your token",
			})
		}

		return c.Next()
	}
}
