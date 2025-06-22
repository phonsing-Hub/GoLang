package jwt

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/phonsing-Hub/GoLang/config"
	"time"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// var SecretKey = []byte(config.Env.JWTSecret)

func getSecretKey() ([]byte, error) {
    if config.Env.JWTSecret == "" {
        // This log should ideally be caught by config.LoadEnv() itself,
        // but it's a safeguard if jwt functions are called too early.
        return nil, fmt.Errorf("JWT_SECRET is not loaded in configuration")
    }
    return []byte(config.Env.JWTSecret), nil
}

func GenerateToken(userID uint, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	SecretKey, _ := getSecretKey()
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "your-api-issuer",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	SecretKey, _ := getSecretKey()
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return SecretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	if c, ok := token.Claims.(*Claims); ok {
		return c, nil
	}
	return nil, fmt.Errorf("invalid token claims type")
}

func GetClaimsFromFiberContext(c *fiber.Ctx) (*Claims, error) {
	claims, ok := c.Locals("user").(*Claims)
	if !ok || claims == nil {
		return nil, fmt.Errorf("claims not found in Fiber context")
	}
	return claims, nil
}
