package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// GenerateJWTWithExpiry creates a JWT containing the user's ID, username,
// role, issuer and audience. The token expires after the given duration.
func GenerateJWTWithExpiry(secret string, userID uint, username, role, issuer, audience string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"iss":      issuer,
		"aud":      audience,
		"exp":      time.Now().Add(expiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// JWTMiddleware เช็ค Access Token จาก Header
func JWTMiddleware(secret, issuer, audience string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Authorization header format"})
		}
		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}
		claims := token.Claims.(jwt.MapClaims)
		if issuer != "" {
			if iss, ok := claims["iss"].(string); !ok || iss != issuer {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token issuer"})
			}
		}
		if audience != "" {
			if aud, ok := claims["aud"]; ok {
				switch v := aud.(type) {
				case string:
					if v != audience {
						return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token audience"})
					}
				case []interface{}:
					match := false
					for _, a := range v {
						if str, ok := a.(string); ok && str == audience {
							match = true
						}
					}
					if !match {
						return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token audience"})
					}
				default:
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token audience"})
				}
			} else {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token audience"})
			}
		}
		// อ่าน user_id และ role จาก claims ใส่ลง Locals
		if userID, ok := claims["user_id"].(float64); ok {
			c.Locals("user_id", uint(userID))
		}
		if role, ok := claims["role"].(string); ok {
			c.Locals("role", role)
		}
		c.Locals("username", claims["username"].(string))
		return c.Next()
	}
}
