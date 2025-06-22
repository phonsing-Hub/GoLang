package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/phonsing-Hub/GoLang/pkg/jwt"
	"github.com/phonsing-Hub/GoLang/utils/response"
)

func JWTAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Fail(c, "UNAUTHORIZED", "Authorization header is required", fiber.StatusUnauthorized)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return response.Fail(c, "UNAUTHORIZED", "Authorization header must be in Bearer token format", fiber.StatusUnauthorized)
		}

		tokenString := parts[1]

		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			return response.Fail(c, "UNAUTHORIZED", "Invalid or expired token: "+err.Error(), fiber.StatusUnauthorized)
		}

		c.Locals("user", claims)

		return c.Next()
	}
}
