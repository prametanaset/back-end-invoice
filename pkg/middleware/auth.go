package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// GenerateJWTWithExpiry สร้าง JWT ระบุอายุเป็น duration ได้
// GenerateJWTWithExpiry creates a JWT containing the user's ID, username and role
func GenerateJWTWithExpiry(secret string, userID uint, username, role string, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(expiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// JWTMiddleware เช็ค Access Token จาก Header
func JWTMiddleware(secret string) fiber.Handler {
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

// JWTMiddlewareExcept returns JWT middleware that skips the check for the
// provided path prefixes. It is useful when applying authentication globally
// but keeping some routes (e.g. /auth) public.
func JWTMiddlewareExcept(secret string, skip ...string) fiber.Handler {
	base := JWTMiddleware(secret)
	return func(c *fiber.Ctx) error {
		path := c.Path()
		for _, prefix := range skip {
			if strings.HasPrefix(path, prefix) {
				return c.Next()
			}
		}
		return base(c)
	}
}
