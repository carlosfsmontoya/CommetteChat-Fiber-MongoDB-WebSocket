package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Obtener el token del encabezado Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or malformed JWT"})
		}

		// Eliminar el prefijo "Bearer " del token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parsear el token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validar el método de firma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
			}
			// Devolver la clave secreta
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
		}

		// Continuar con la siguiente función en la pila de middleware
		return c.Next()
	}
}

func SecretKeyRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Obtener la clave secreta del encabezado
		secretKey := c.Get("X-Secret-Key")
		if secretKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing secret key"})
		}

		// Comparar la clave secreta con la esperada
		expectedSecretKey := os.Getenv("FAST_GO_KEY")
		if secretKey != expectedSecretKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid secret key"})
		}

		// Continuar con la siguiente función en la pila de middleware
		return c.Next()
	}
}
