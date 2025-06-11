package middleware

import (
	"invoice_project/pkg/apperror"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// GenerateJWTWithExpiry creates a JWT using a common set of OAuth-style claims
// and embeds a token_type so that access and refresh tokens can be
// differentiated.
func GenerateJWTWithExpiry(secret string, userID uint, role string, expiry time.Duration, tokenType string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iss":        "https://auth.example.com",
		"sub":        strconv.FormatUint(uint64(userID), 10),
		"aud":        "https://api.example.com",
		"scope":      role,
		"token_type": tokenType,
		"exp":        now.Add(expiry).Unix(),
		"iat":        now.Unix(),
		"nbf":        now.Unix(),
		"jti":        uuid.NewString(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// JWTMiddleware เช็ค Access Token จาก Header
func JWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return apperror.New(fiber.StatusUnauthorized)
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
		if typ, ok := claims["token_type"].(string); !ok || typ != "access" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token type"})
		}
		// Map standard claims back to our context locals
		if sub, ok := claims["sub"].(string); ok {
			if id, err := strconv.Atoi(sub); err == nil {
				c.Locals("user_id", uint(id))
			}
		}
		if scope, ok := claims["scope"].(string); ok {
			c.Locals("role", scope)
		}
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
